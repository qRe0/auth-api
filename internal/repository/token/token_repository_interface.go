package token

import "context"

type TokenRepositryInterface interface {
	SaveToken(ctx context.Context, token, userID string) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
	CheckRefreshToken(ctx context.Context, key string) (string, error)
}
