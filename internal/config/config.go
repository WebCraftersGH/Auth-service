package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPPort       string
	LoggingLevel   string
	JWTSecret      string
	JWTTTL         time.Duration
	OTPLength      int
	OTPTTL         time.Duration
	OTPMaxAttempts int

	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	KafkaBrokers []string
	KafkaTopic   string
	KafkaTimeout time.Duration
	KafkaGroupID string
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load() (Config, error) {
	cfg := Config{
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		LoggingLevel:   strings.ToUpper(getEnv("LOGGING_LEVEL", "INFO")),
		JWTSecret:      getEnv("JWT_SECRET", "super-secret-change-me"),
		JWTTTL:         getEnvDuration("JWT_TTL", 24*time.Hour),
		OTPLength:      getEnvInt("OTP_LENGTH", 6),
		OTPTTL:         getEnvDuration("OTP_TTL", 10*time.Minute),
		OTPMaxAttempts: getEnvInt("OTP_MAX_ATTEMPTS", 5),

		DBUser: getEnv("DB_USER", "authuser"),
		DBPass: getEnv("DB_PASS", "authpass"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBName: getEnv("DB_NAME", "authdb"),

		SMTPHost: getEnv("SMTP_HOST", "localhost"),
		SMTPPort: getEnvInt("SMTP_PORT", 1025),
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASS", ""),
		SMTPFrom: getEnv("SMTP_FROM", "no-reply@example.com"),

		KafkaBrokers: splitCSV(getEnv("KAFKA_BROKERS", "localhost:9092")),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "users.create.requested"),
		KafkaTimeout: getEnvDuration("KAFKA_TIMEOUT", 5*time.Second),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "auth-service"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.OTPLength <= 0 {
		return Config{}, fmt.Errorf("OTP_LENGTH must be > 0")
	}
	if cfg.OTPMaxAttempts <= 0 {
		return Config{}, fmt.Errorf("OTP_MAX_ATTEMPTS must be > 0")
	}
	if len(cfg.KafkaBrokers) == 0 {
		return Config{}, fmt.Errorf("KAFKA_BROKERS is required")
	}

	return cfg, nil
}

func (c Config) DBConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return value
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
