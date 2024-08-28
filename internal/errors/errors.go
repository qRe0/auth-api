package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrLoadEnvVars   = errors.New("failed to load environment variables")
	ErrIncorrectData = errors.New("incorrect data")
	ErrCreateUser    = errors.New("failed to create user")
	ErrUserNotFound  = errors.New("user not found in database")
	ErrGenerateToken = errors.New("failed to generate token")
	ErrSaveToken     = errors.New("failed to save token")
	ErrTokenNotFound = errors.New("token not found")
	ErrDeletingToken = errors.New("failed to delete token")
)
