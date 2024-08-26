package auth

import (
	"context"

	"github.com/qRe0/auth-api/internal/models"
)

type AuthServiceInterface interface {
	SignUp(ctx context.Context, user *models.User) error
}
