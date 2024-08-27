package auth

import (
	"context"

	"github.com/qRe0/auth-api/internal/models"
)

type AuthRepositoryInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserData(ctx context.Context, phone string) (models.User, error)
}
