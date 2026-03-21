package tokensrepo

import (
	"time"

	"github.com/google/uuid"
)

type tokenModel struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Value     string    `gorm:"column:value"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;index"`
}

func (tokenModel) TableName() string {
	return "tokens"
}

type tokenRow struct {
	ID        uuid.UUID
	Value     string
	UserID    uuid.UUID
	Email     string
	CreatedAt time.Time
	ExpiredAt time.Time
}
