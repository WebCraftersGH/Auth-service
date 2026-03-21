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
	"github.com/google/uuid"
)

type authSVC struct {
	usersRepo     contracts.UsersRepo
	mailSVC       contracts.MailSVC
	otpSVC        contracts.OTPSVC
	tokenSVC      contracts.TokenSVC
	eventsProduce contracts.UserEventsProducer
	logger        contracts.ILogger
}

func NewAuthSVC(
	usersRepo contracts.UsersRepo,
	mailSVC contracts.MailSVC,
	otpSVC contracts.OTPSVC,
	tokenSVC contracts.TokenSVC,
	eventsProducer contracts.UserEventsProducer,
	logger contracts.ILogger,
) *authSVC {
	return &authSVC{
		usersRepo:     usersRepo,
		mailSVC:       mailSVC,
		otpSVC:        otpSVC,
		tokenSVC:      tokenSVC,
		eventsProduce: eventsProducer,
		logger:        logger,
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

	s.logger.Debugf("start auth: otp generated and sent, email=%s", e)
	return nil
}

func (s *authSVC) OTPCheck(ctx context.Context, email, code string) (domain.Token, error) {
	e, err := s.validateEmail(email)
	if err != nil {
		s.logger.Debugf("otp check: invalid email: %s", email)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrInvalidEmail, err)
	}

	if err = s.otpSVC.VerifyCode(ctx, e, code); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOTP):
			return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrInvalidOTP, err)
		case errors.Is(err, domain.ErrToManyOTPAttempts):
			return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrToManyOTPAttempts, err)
		case errors.Is(err, domain.ErrOTPExpired):
			return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrOTPExpired, err)
		default:
			return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
		}
	}

	u, createdNow, err := s.ensureUser(ctx, e)
	if err != nil {
		s.logger.Errorf("otp check: ensure user failed, email=%s err=%v", e, err)
		return domain.Token{}, err
	}

	if createdNow {
		event := domain.UserCreateRequestedEvent{
			EventID:   uuid.New(),
			UserID:    u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}
		if err := s.eventsProduce.PublishUserCreateRequested(ctx, event); err != nil {
			_ = s.usersRepo.DeleteUserByID(ctx, u.ID)
			s.logger.Errorf("otp check: failed to publish user-create event, email=%s err=%v", e, err)
			return domain.Token{}, fmt.Errorf("%w: %v", domain.ErrKafkaPublish, err)
		}
	}

	t, err := s.tokenSVC.ReadToken(ctx, u.Email)
	if err == nil {
		return t, nil
	}
	if !errors.Is(err, domain.ErrTokenNotFound) && !errors.Is(err, domain.ErrTokenExpired) {
		s.logger.Errorf("otp check: read token failed, email=%s err=%v", e, err)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	t, err = s.tokenSVC.GenerateJWT(u)
	if err != nil {
		s.logger.Errorf("otp check: generate jwt failed, email=%s err=%v", e, err)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	if err := s.tokenSVC.SaveToken(ctx, t); err != nil {
		s.logger.Errorf("otp check: save token failed, email=%s err=%v", e, err)
		return domain.Token{}, fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	return t, nil
}

func (s *authSVC) AuthCheck(ctx context.Context, token string) error {
	claims, err := s.tokenSVC.ParseToken(token)
	if err != nil {
		return err
	}

	stored, err := s.tokenSVC.ReadToken(ctx, claims.UserEmail)
	if err != nil {
		if errors.Is(err, domain.ErrTokenExpired) || errors.Is(err, domain.ErrTokenNotFound) {
			return domain.ErrUnauthorized
		}
		return fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	if stored.Value != token {
		return domain.ErrUnauthorized
	}

	return nil
}

func (s *authSVC) Logout(ctx context.Context, email string) error {
	e, err := s.validateEmail(email)
	if err != nil {
		s.logger.Debugf("logout: invalid email: %s", email)
		return fmt.Errorf("%w: %v", domain.ErrInvalidEmail, err)
	}

	err = s.tokenSVC.DeleteToken(ctx, e)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	return nil
}

func (s *authSVC) ensureUser(ctx context.Context, email string) (domain.User, bool, error) {
	u, err := s.usersRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return u, false, nil
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return domain.User{}, false, fmt.Errorf("%w: %v", domain.InternalError, err)
	}

	created, createErr := s.usersRepo.CreateUser(ctx, domain.User{Email: email})
	if createErr != nil {
		return domain.User{}, false, fmt.Errorf("%w: %v", domain.InternalError, createErr)
	}

	return created, true, nil
}

func (s *authSVC) validateEmail(raw string) (string, error) {
	str := strings.TrimSpace(raw)
	if str == "" {
		return "", fmt.Errorf("%w: empty email", domain.ErrInvalidEmail)
	}

	if utf8.RuneCountInString(str) > 254 {
		return "", fmt.Errorf("%w: email > 254 chars", domain.ErrInvalidEmail)
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
