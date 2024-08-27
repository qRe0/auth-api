package token

import "context"

type TokenRepositryInterface interface {
	SaveToken(ctx context.Context, token, userID string) error
}
