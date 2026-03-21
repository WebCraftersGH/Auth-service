package usersrepo

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userModel struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string         `gorm:"column:email"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (userModel) TableName() string {
	return "users"
}
