package token

import (
	"context"

	"github.com/qRe0/auth-api/internal/models"
)

type TokenServiceInterface interface {
	SaveToken(ctx context.Context, token, userID string) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
	RefreshTokenExists(ctx context.Context, userID string) (models.RefreshTokenExistsResponse, error)
}
