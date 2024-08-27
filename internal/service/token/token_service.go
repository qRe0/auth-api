package token

import (
	"context"

	errs "github.com/qRe0/auth-api/internal/errors"
	repository "github.com/qRe0/auth-api/internal/repository/token"
)

type TokenService struct {
	repo repository.TokenRepositryInterface
}

func NewTokenService(repo repository.TokenRepositryInterface) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (t *TokenService) SaveToken(ctx context.Context, token string, userID string) error {
	err := t.repo.SaveToken(ctx, token, userID)
	if err != nil {
		return errs.ErrSaveToken
	}

	return nil
}
