package tokensrepo

import (
	"context"
	"errors"

	"github.com/WebCraftersGH/Auth-service/internal/domain"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) ReadTokenByEmail(ctx context.Context, email string) (domain.Token, error) {
	var row tokenRow
	err := r.db.WithContext(ctx).
		Table("tokens AS t").
		Select("t.id, t.value, t.user_id, u.email, t.created_at, t.expired_at").
		Joins("JOIN users u ON u.id = t.user_id").
		Where("u.email = ? AND u.deleted_at IS NULL", email).
		Order("t.created_at DESC").
		Limit(1).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Token{}, domain.ErrTokenNotFound
		}
		return domain.Token{}, err
	}

	return toDomain(row)
}

func (r *Repo) DeleteTokenByEmail(ctx context.Context, email string) error {
	subQuery := r.db.WithContext(ctx).
		Table("users").
		Select("id").
		Where("email = ? AND deleted_at IS NULL", email)

	res := r.db.WithContext(ctx).
		Where("user_id IN (?)", subQuery).
		Delete(&tokenModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrTokenNotFound
	}

	return nil
}

func (r *Repo) SaveToken(ctx context.Context, t domain.Token) error {
	model := tokenModel{
		Value:     t.Value,
		ExpiredAt: t.ExpiredAt,
		UserID:    t.User.ID,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}
