package usecase

import (
	"github.com/WebCraftersGH/Auth-service/internal/contracts"
)

type authSVC struct {
	userRepo contracts.UsersRepo
	mailSVC  contracts.MailSVC
	otpSVC   contracts.OTPSVC
	tokenSVC contracts.TokenSVC
}

func NewAuthSVC() *authSVC {
	return &authSVC{}
}
