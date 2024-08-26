package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/qRe0/auth-api/configs"
	errs "github.com/qRe0/auth-api/internal/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/qRe0/auth-api/internal/models"
	repository "github.com/qRe0/auth-api/internal/repository/auth"
)

type AuthService struct {
	cfg  configs.JWTConfig
	repo repository.AuthRepositoryInterface
}

func NewAuthService(cfg configs.JWTConfig, repo repository.AuthRepositoryInterface) *AuthService {
	return &AuthService{
		cfg:  cfg,
		repo: repo,
	}
}

func (a *AuthService) SignUp(ctx context.Context, user *models.User) error {
	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		return errors.Wrap(errs.ErrIncorrectData, "failed to validate registration data. incorrect user data")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	user.Password = hashedPassword
	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, errs.ErrCreateUser.Error())
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}
