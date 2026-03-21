package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WebCraftersGH/Auth-service/internal/config"
	"github.com/WebCraftersGH/Auth-service/internal/controller"
	appkafka "github.com/WebCraftersGH/Auth-service/internal/kafka"
	mailsrepo "github.com/WebCraftersGH/Auth-service/internal/repository/mails_repo"
	otpsrepo "github.com/WebCraftersGH/Auth-service/internal/repository/otps_repo"
	tokensrepo "github.com/WebCraftersGH/Auth-service/internal/repository/tokens_repo"
	usersrepo "github.com/WebCraftersGH/Auth-service/internal/repository/users_repo"
	"github.com/WebCraftersGH/Auth-service/internal/usecase"
	"github.com/WebCraftersGH/Auth-service/pkg/logging"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()
	logging.Init(cfg.LoggingLevel)
	logger := logging.GetLogger()

	db, err := gorm.Open(postgres.Open(cfg.DBConnString()), &gorm.Config{})
	if err != nil {
		logger.Fatalf("open db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("get sql db: %v", err)
	}
	defer sqlDB.Close()

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err = sqlDB.Ping(); err != nil {
		logger.Fatalf("ping db: %v", err)
	}

	otpRepo := otpsrepo.New()
	usersRepo := usersrepo.New(db)
	tokensRepo := tokensrepo.New(db)
	mailsRepo := mailsrepo.New(db)
	producer := appkafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaTimeout, logger)
	defer producer.Close()

	otpSVC := usecase.NewOTPSVC(otpRepo, cfg.OTPLength, cfg.OTPTTL, cfg.OTPMaxAttempts, logger)
	mailSVC := usecase.NewMailSVC(mailsRepo, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom, logger)
	tokenSVC := usecase.NewTokenSVC(tokensRepo, cfg.JWTSecret, cfg.JWTTTL, logger)
	authSVC := usecase.NewAuthSVC(usersRepo, mailSVC, otpSVC, tokenSVC, producer, logger)

	mux := http.NewServeMux()
	controller.New(authSVC, logger).Register(mux)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Infof("auth-service is listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("server shutdown: %v", err)
	}
}
