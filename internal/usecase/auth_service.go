package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/WebCraftersGH/Auth-service/internal/contracts"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
)

type authSVC struct {
	usersRepo contracts.UsersRepo
	mailSVC   contracts.MailSVC
	otpSVC    contracts.OTPSVC
	tokenSVC  contracts.TokenSVC
	logger    contracts.ILogger
}

func NewAuthSVC(
	usersRepo contracts.UsersRepo,
	mailSVC contracts.MailSVC,
	otpSVC contracts.OTPSVC,
	tokenSVC contracts.TokenSVC,
	logger contracts.ILogger,
) *authSVC {
	return &authSVC{
		usersRepo: usersRepo,
		mailSVC:   mailSVC,
		otpSVC:    otpSVC,
		tokenSVC:  tokenSVC,
		logger:    logger,
	}
}

func (s *authSVC) StartAuth(ctx context.Context, email string) error {

	e, err := s.validateEmail(email)
	if err != nil {
		s.logger.Debugf("start auth: invalid email: %s", email)
		return fmt.Errorf("%w: %v", domain.ErrInvalidEmail, err)
	}

	code := s.otpSVC.GenerateOTP()
	code.UserEmail = e

	if err := s.otpSVC.SaveCode(ctx, code); err != nil {
		s.logger.Errorf("start auth: failed to save otp, email=%s, err=%s", e, err)
		return fmt.Errorf("%w: %v", domain.ErrOtpSave, err)
	}

	if err := s.mailSVC.SendMail(ctx, e, code); err != nil {
		s.logger.Errorf("start auth: failed to send mail, email=%s, err=%s", e, err)
		return fmt.Errorf("%w: %v", domain.ErrSendMail, err)
	}

	s.logger.Debugf("start auth: otp generated and send, email=%s", e)
	return nil
}

func (s *authSVC) OTPCheck(ctx context.Context, email, code string) (domain.Token, error) {

	e, err := s.validateEmail(email)
	if err != nil {
		s.logger.Debugf("otp check: invalid email: %s", email)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrInvalidEmail, err)
	}

	err = s.otpSVC.VerifyCode(ctx, e, code)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidOTP) {
			s.logger.Debugf("otp check: invalid otp, email=%s", e)
			return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrInvalidOTP, err)
		} else if errors.Is(err, domain.ErrToManyOTPAttempts) {
			s.logger.Debugf("otp check: to many attempts, email=%s, otp=%s", e, code)
			return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrToManyOTPAttempts, err)
		}
		s.logger.Errorf("otp check: failed to verify otp code, email=%s, otp=%s", e, code)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	var u domain.User

	u, err = s.usersRepo.GetUserByEmail(ctx, e)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// If user not found.
			user := domain.User{
				Email: e,
			}

			u, err = s.usersRepo.CreateUser(ctx, user)
			if err != nil {
				s.logger.Errorf("otp check: failed to create new user, email=%s", e)
				return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
			}
		}
		s.logger.Errorf("otp check: get user failed, email=%s", e)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	var t domain.Token
	t, err = s.tokenSVC.ReadToken(ctx, u.Email)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) || errors.Is(err, domain.ErrTokenExpired) {
			//Generate and save token
			t, err = s.tokenSVC.GenerateJWT(u)
			if err != nil {
				s.logger.Errorf("otp check: generate jwt failed, email=%s", e)
				return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
			}
		}
		s.logger.Errorf("otp check: read token failed, email=%s", e)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	return t, nil
}

func (s *authSVC) AuthCheck(ctx context.Context, token string) error

func (s *authSVC) Logout(ctx context.Context, email string) error {

	e, err := s.validateEmail(email)
	if err != nil {
		s.logger.Debugf("logout: invalid email: %s", email)
		return fmt.Errorf("%w: %v", domain.ErrInvalidEmail, err)
	}

	err = s.tokenSVC.DeleteToken(ctx, e)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			s.logger.Debugf("logout: user already logout, email=%s", e)
			return err
		}
		s.logger.Errorf("logout: delete token error, email=%s", e)
		return fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	return nil
}

func (s *authSVC) validateEmail(raw string) (string, error) {
	str := strings.TrimSpace(raw)
	if str == "" {
		return "", fmt.Errorf("%w: Empty email!", domain.ErrInvalidEmail)
	}

	if utf8.RuneCountInString(str) > 254 {
		return "", fmt.Errorf("%w: Email > 254 chars", domain.ErrInvalidEmail)
	}

	addr, err := mail.ParseAddress(str)
	if err != nil {
		return "", domain.ErrInvalidEmail
	}
	if addr.Address != str {
		return "", domain.ErrInvalidEmail
	}

	return strings.ToLower(str), nil
}
