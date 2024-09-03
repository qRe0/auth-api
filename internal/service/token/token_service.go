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

func (t *TokenService) SaveToken(ctx context.Context, token string, userID int) error {
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

func (t *TokenService) RefreshTokenExists(ctx context.Context, userID int) (models.RefreshTokenExistsResponse, error) {
	key := fmt.Sprintf("user:%d", userID)
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

func (t *TokenService) GetToken(ctx context.Context, userID int) (string, error) {
	key := fmt.Sprintf("user:%d", userID)
	token, err := t.repo.GetToken(ctx, key)
	if err != nil {
		return "", errs.ErrTokenNotFound
	}

	return token, nil
}

func (t *TokenService) DeleteToken(ctx context.Context, userID int) error {
	key := fmt.Sprintf("user:%d", userID)
	err := t.repo.DeleteToken(ctx, key)
	if err != nil {
		return errs.ErrDeletingToken
	}

	return nil
}

func (t *TokenService) TokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklisted:%s", token)

	return t.repo.TokenBlacklisted(ctx, key)
}

func (t *TokenService) BlacklistToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("blacklisted:%s", token)
	err := t.repo.BlacklistToken(ctx, key)
	if err != nil {
		return errs.ErrBlacklistingToken
	}

	return nil
}
