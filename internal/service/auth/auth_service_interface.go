package auth

import (
	"context"
	"time"

	"github.com/qRe0/auth-api/internal/models"
)

type AuthServiceInterface interface {
	SignUp(ctx context.Context, user *models.User) (models.Tokens, error)
	LogIn(ctx context.Context, user *models.User) (models.Tokens, error)
	NewSession(ctx context.Context, userID int, secretKey string, lifetime time.Duration) (models.Tokens, error)
	Refresh(ctx context.Context, token string) (models.Tokens, error)
	RevokeTokens(ctx context.Context, user *models.User) error
	ValidateToken(token string, cfg string) (string, error)
	TokenBlacklisted(ctx context.Context, token string) (bool, error)
}
