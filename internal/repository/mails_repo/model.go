package mailsrepo

import (
	"time"

	"github.com/google/uuid"
)

type mailModel struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Value     string    `gorm:"column:value"`
	ToEmail   string    `gorm:"column:to_email"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (mailModel) TableName() string {
	return "mails"
}
