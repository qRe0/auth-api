package token

import "context"

type TokenRepositryInterface interface {
	SaveToken(ctx context.Context, token string, userID int) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
	CheckRefreshToken(ctx context.Context, key string) (string, error)
	GetToken(ctx context.Context, key string) (string, error)
	DeleteToken(ctx context.Context, key string) error
	TokenBlacklisted(ctx context.Context, key string) (bool, error)
}
