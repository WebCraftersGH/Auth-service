package mailsrepo

import (
	"context"
	"errors"

	"github.com/WebCraftersGH/Auth-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) SaveMail(ctx context.Context, m domain.Mail) error {
	model := mailModel{
		Value:   m.Value,
		ToEmail: m.ToEmail,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *Repo) GetMail(ctx context.Context, mailID uuid.UUID) (domain.Mail, error) {
	var model mailModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", mailID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Mail{}, domain.InternalError
		}
		return domain.Mail{}, err
	}

	return toDomain(model)
}

func (r *Repo) DeleteMail(ctx context.Context, mailID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&mailModel{}, "id = ?", mailID).Error
}
