package usecase

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"github.com/WebCraftersGH/Auth-service/internal/contracts"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
)

const (
	LETTERS        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DEFAULT_LENGTH = 6
)

type otpSVC struct {
	otpRepo        contracts.OTPsRepo
	otpLength      int
	otpTTL         time.Duration
	otpMaxAttempts int
	logger         contracts.ILogger
}

func NewOTPSVC(
	otpRepo contracts.OTPsRepo,
	otpLength int,
	otpTTL time.Duration,
	otpMaxAttempts int,
	logger contracts.ILogger,
) *otpSVC {
	if otpLength <= 0 {
		otpLength = DEFAULT_LENGTH
	}
	if otpMaxAttempts <= 0 {
		otpMaxAttempts = 5
	}
	if otpTTL <= 0 {
		otpTTL = 10 * time.Minute
	}

	return &otpSVC{
		otpRepo:        otpRepo,
		otpLength:      otpLength,
		otpTTL:         otpTTL,
		otpMaxAttempts: otpMaxAttempts,
		logger:         logger,
	}
}

func (s *otpSVC) GenerateOTP() domain.OTP {
	code := make([]byte, s.otpLength)
	for i := 0; i < s.otpLength; i++ {
		randInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))
		code[i] = LETTERS[randInt.Int64()]
	}

	return domain.OTP{
		Value:     string(code),
		Attempts:  0,
		ExpiresAt: time.Now().Add(s.otpTTL),
	}
}

func (s *otpSVC) SaveCode(_ context.Context, code domain.OTP) error {
	return s.otpRepo.SaveCode(code, strings.ToLower(strings.TrimSpace(code.UserEmail)))
}

func (s *otpSVC) VerifyCode(_ context.Context, email, code string) error {
	otp, err := s.otpRepo.GetCode(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return err
	}

	if time.Now().After(otp.ExpiresAt) {
		_ = s.otpRepo.DeleteCode(email)
		return domain.ErrOTPExpired
	}

	if otp.Attempts >= s.otpMaxAttempts {
		_ = s.otpRepo.DeleteCode(email)
		return domain.ErrToManyOTPAttempts
	}

	if !strings.EqualFold(strings.TrimSpace(otp.Value), strings.TrimSpace(code)) {
		if err := s.otpRepo.IncAttempts(email); err != nil {
			return err
		}

		refreshed, err := s.otpRepo.GetCode(email)
		if err != nil {
			return err
		}
		if refreshed.Attempts >= s.otpMaxAttempts {
			_ = s.otpRepo.DeleteCode(email)
			return domain.ErrToManyOTPAttempts
		}
		return domain.ErrInvalidOTP
	}

	return s.otpRepo.DeleteCode(email)
}
