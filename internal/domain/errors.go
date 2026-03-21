package domain

import "errors"

var (
	ErrInvalidEmail = errors.New("invalid email")

	ErrOtpSave           = errors.New("otp save failed")
	ErrInvalidOTP        = errors.New("invalid otp code")
	ErrToManyOTPAttempts = errors.New("too many attempts")
	ErrOTPExpired        = errors.New("otp expired")
	ErrOTPNotFound       = errors.New("otp not found")

	ErrSendMail = errors.New("send mail failed")

	ErrUserNotFound = errors.New("user not found")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
	ErrUnauthorized  = errors.New("unauthorized")

	ErrKafkaPublish = errors.New("kafka publish failed")

	InternalError = errors.New("internal error")
)
