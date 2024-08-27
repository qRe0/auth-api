package token

import "context"

type TokenRepositryInterface interface {
	SaveToken(ctx context.Context, token string, userID int) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
	CheckRefreshToken(ctx context.Context, key string) (string, error)
}
