package usecase

import (
	"context"
	"crypto/rand"
	"github.com/WebCraftersGH/Auth-service/internal/contracts"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
	"math/big"
)

const (
	LETTERS        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DEFAULT_LENGTH = 6
)

type otpSVC struct {
	otpRepo   contracts.OTPsRepo
	otpLength int
	logger    contracts.ILogger
}

func NewOTPSVC(
	otpRepo contracts.OTPsRepo,
	otpLength int,
	logger contracts.ILogger,
) *otpSVC {

	if otpLength <= 0 || otpLength >= len(LETTERS) {
		otpLength = DEFAULT_LENGTH
	}

	return &otpSVC{
		otpRepo:   otpRepo,
		otpLength: otpLength,
		logger:    logger,
	}
}

func (s *otpSVC) GenerateOTP() domain.OTP {

	code := make([]byte, s.otpLength)
	for i := 0; i < s.otpLength; i++ {
		randInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))
		code[i] = LETTERS[randInt.Int64()]
	}

	return domain.OTP{Value: string(code)}
}

func (s *otpSVC) SaveCode(ctx context.Context, code domain.OTP) error {
	return nil
}

func (s *otpSVC) VerifyCode(ctx context.Context, email, code string) error {
	return nil
}
