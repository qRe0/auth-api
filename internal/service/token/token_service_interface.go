package token

import "context"

type TokenServiceInterface interface {
	SaveToken(ctx context.Context, token, userID string) error
}
