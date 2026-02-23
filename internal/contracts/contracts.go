package contracts

import (
	"context"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
	"github.com/google/uuid"
)

type AuthSVC interface {
	StartAuth(ctx context.Context, email string) error
	OTPCheck(ctx context.Context, email, code string) (domain.Token, error)
	AuthCheck(ctx context.Context, token string) error
	Logout(ctx context.Context, token string) error
}

type TokenSVC interface {
	GenerateJWT(user domain.User) (domain.Token, error)
	SaveToken(ctx context.Context, token domain.Token) error
	ReadToken(ctx context.Context, email string) (domain.Token, error)
	DeleteToken(ctx context.Context, email string) error
}

type MailSVC interface {
	SendMail(ctx context.Context, toEmail, code domain.OTP) error
}

type OTPSVC interface {
	GenerateOTP() domain.OTP
	SaveCode(ctx context.Context, code domain.OTP) error
	VerifyCode(ctx context.Context, email, code string) error
}

type UsersRepo interface {
	SaveUser(ctx context.Context, u domain.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
}

type OTPsRepo interface {
	SaveCode(otp domain.OTP, key string) error
	GetCode(key string) (domain.OTP, error)
	IncAttempts(key string) error
	DeleteCode(key string) error
}

type TokensRepo interface {
	ReadTokenByEmail(ctx context.Context, email string) (domain.Token, error)
	DeleteTokenByEmail(ctx context.Context, email string) error
	SaveToken(ctx context.Context, t domain.Token) error
}

type MailsRepo interface {
	SaveMail(ctx context.Context, m domain.Mail) error
	GetMail(ctx context.Context, mailID uuid.UUID) (domain.Mail, error)
	DeleteMail(ctx context.Context, mailID uuid.UUID) error
}
