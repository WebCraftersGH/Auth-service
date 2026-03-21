package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/WebCraftersGH/Auth-service/internal/contracts"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenSVC struct {
	tokensRepo contracts.TokensRepo
	jwtSecret  []byte
	jwtTTL     time.Duration
	logger     contracts.ILogger
}

type jwtClaims struct {
	UserID    string `json:"uid"`
	UserEmail string `json:"email"`
	jwt.RegisteredClaims
}

func NewTokenSVC(tokensRepo contracts.TokensRepo, jwtSecret string, jwtTTL time.Duration, logger contracts.ILogger) *tokenSVC {
	if jwtTTL <= 0 {
		jwtTTL = 24 * time.Hour
	}

	return &tokenSVC{
		tokensRepo: tokensRepo,
		jwtSecret:  []byte(jwtSecret),
		jwtTTL:     jwtTTL,
		logger:     logger,
	}
}

func (s *tokenSVC) GenerateJWT(user domain.User) (domain.Token, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtTTL)
	tokenID := uuid.New()

	claims := jwtClaims{
		UserID:    user.ID.String(),
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return domain.Token{}, err
	}

	return domain.Token{
		ID:        tokenID,
		Value:     signed,
		User:      user,
		CreatedAt: now,
		ExpiredAt: expiresAt,
	}, nil
}

func (s *tokenSVC) SaveToken(ctx context.Context, token domain.Token) error {
	return s.tokensRepo.SaveToken(ctx, token)
}

func (s *tokenSVC) ReadToken(ctx context.Context, email string) (domain.Token, error) {
	token, err := s.tokensRepo.ReadTokenByEmail(ctx, email)
	if err != nil {
		return domain.Token{}, err
	}

	if time.Now().After(token.ExpiredAt) {
		_ = s.tokensRepo.DeleteTokenByEmail(ctx, email)
		return domain.Token{}, domain.ErrTokenExpired
	}

	return token, nil
}

func (s *tokenSVC) DeleteToken(ctx context.Context, email string) error {
	return s.tokensRepo.DeleteTokenByEmail(ctx, email)
}

func (s *tokenSVC) ParseToken(raw string) (domain.TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(raw, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return domain.TokenClaims{}, domain.ErrTokenExpired
		}
		return domain.TokenClaims{}, domain.ErrUnauthorized
	}

	claims, ok := parsed.Claims.(*jwtClaims)
	if !ok || !parsed.Valid {
		return domain.TokenClaims{}, domain.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return domain.TokenClaims{}, domain.ErrUnauthorized
	}

	if claims.ExpiresAt == nil {
		return domain.TokenClaims{}, domain.ErrUnauthorized
	}

	return domain.TokenClaims{
		UserID:    userID,
		UserEmail: claims.UserEmail,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
