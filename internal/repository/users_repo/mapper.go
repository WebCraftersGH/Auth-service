package usersrepo

import "github.com/WebCraftersGH/Auth-service/internal/domain"

func toDomain(model userModel) (domain.User, error) {
	return domain.User{
		ID:        model.ID,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
	}, nil
}
