package token

import "context"

type TokenServiceInterface interface {
	SaveToken(ctx context.Context, token, userID string) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (string, error)
}
