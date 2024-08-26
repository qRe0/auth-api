package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrLoadEnvVars   = errors.New("failed to load environment variables")
	ErrIncorrectData = errors.New("incorrect data")
	ErrCreateUser    = errors.New("failed to create user")
)
