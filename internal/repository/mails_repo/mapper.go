package mailsrepo

import "github.com/WebCraftersGH/Auth-service/internal/domain"

func toDomain(model mailModel) (domain.Mail, error) {
	return domain.Mail{
		ID:        model.ID,
		Value:     model.Value,
		ToEmail:   model.ToEmail,
		CreatedAt: model.CreatedAt,
	}, nil
}
