package token

import (
	"context"

	"github.com/pkg/errors"
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

func (t *TokenService) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
	userID, err := t.repo.GetUserIDByRefreshToken(ctx, token)
	if err != nil {
		return "", errors.Wrap(err, "failed to get user id by refresh token from Cache")
	}

	return userID, nil
}
