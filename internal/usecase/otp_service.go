package usecase

import (
	"context"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
)

type otpSVC struct{}

func NewOTPSVC() *otpSVC {
	return &otpSVC{}
}

func (s *otpSVC) GenerateOTP() domain.OTP {
	return domain.OTP{}
}

func (s *otpSVC) SaveCode(ctx context.Context, code domain.OTP) error {
	return nil
}

func (s *otpSVC) VerifyCode(ctx context.Context, email, code string) error {
	return nil
}
