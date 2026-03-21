package usersrepo

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

func (r *Repo) CreateUser(ctx context.Context, u domain.User) (domain.User, error) {
	model := userModel{Email: u.Email}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.User{}, err
	}

	return toDomain(model)
}

func (r *Repo) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	var model userModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return toDomain(model)
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var model userModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return toDomain(model)
}

func (r *Repo) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&userModel{}, "id = ?", userID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
