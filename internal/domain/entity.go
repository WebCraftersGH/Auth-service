package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
}

type Mail struct {
	ID        uuid.UUID
	Value     string
	ToEmail   string
	CreatedAt time.Time
}

type OTP struct {
	Value     string
	UserEmail string
	Attempts  int
}

type Token struct {
	ID        uuid.UUID
	Value     string
	User      User
	CreatedAt time.Time
	ExpiredAt time.Time
}
