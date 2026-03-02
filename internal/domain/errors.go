package domain

import (
	"errors"
)

var (
	ErrInvalidEmail = errors.New("Invalid email")

	ErrOtpSave           = errors.New("Otp save failed")
	ErrInvalidOTP        = errors.New("Invalid otp code")
	ErrToManyOTPAttempts = errors.New("To many attempts")

	ErrSendMail = errors.New("Send mail failed")

	ErrUserNotFound = errors.New("User not found")

	ErrTokenNotFound = errors.New("Token not found")
	ErrTokenExpired  = errors.New("Token expired")
	ErrUnauthorized  = errors.New("Unauthorized")

	InternalError = errors.New("Internal error")
)
