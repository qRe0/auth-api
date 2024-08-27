package token

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	errs "github.com/qRe0/auth-api/internal/errors"
	"github.com/qRe0/auth-api/internal/models"
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

func (t *TokenService) RefreshTokenExists(ctx context.Context, userID string) (models.RefreshTokenExistsResponse, error) {
	key := fmt.Sprintf("user:%s", userID)
	token, err := t.repo.CheckRefreshToken(ctx, key)
	if err != nil {
		return models.RefreshTokenExistsResponse{
			Exists:       false,
			RefreshToken: "not found",
		}, errs.ErrTokenNotFound
	}

	resp := models.RefreshTokenExistsResponse{
		Exists:       true,
		RefreshToken: token,
	}
	return resp, nil
}
