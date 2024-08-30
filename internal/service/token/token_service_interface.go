package token

import (
	"context"

	"github.com/qRe0/auth-api/internal/models"
)

type TokenServiceInterface interface {
	SaveToken(ctx context.Context, token string, userID int) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
	RefreshTokenExists(ctx context.Context, userID int) (models.RefreshTokenExistsResponse, error)
	GetToken(ctx context.Context, userID int) (string, error)
	DeleteToken(ctx context.Context, id int) error
	TokenBlacklisted(ctx context.Context, token string) (bool, error)
}
