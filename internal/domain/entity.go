package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Mail struct {
	ID        uuid.UUID `json:"id"`
	Value     string    `json:"value"`
	ToEmail   string    `json:"to_email"`
	CreatedAt time.Time `json:"created_at"`
}

type OTP struct {
	Value     string    `json:"value"`
	UserEmail string    `json:"user_email"`
	Attempts  int       `json:"attempts"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Token struct {
	ID        uuid.UUID `json:"id"`
	Value     string    `json:"value"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type TokenClaims struct {
	UserID    uuid.UUID
	UserEmail string
	ExpiresAt time.Time
}

type UserCreateRequestedEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
