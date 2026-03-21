package tokensrepo

import "github.com/WebCraftersGH/Auth-service/internal/domain"

func toDomain(model tokenRow) (domain.Token, error) {
	return domain.Token{
		ID:    model.ID,
		Value: model.Value,
		User: domain.User{
			ID:    model.UserID,
			Email: model.Email,
		},
		CreatedAt: model.CreatedAt,
		ExpiredAt: model.ExpiredAt,
	}, nil
}
