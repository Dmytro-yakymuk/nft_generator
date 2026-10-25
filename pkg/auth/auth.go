package auth

import (
	"context"

	"github.com/zeebo/errs"
)

// Error is an error class for auth errors.
var Error = errs.Class("auth error")

// Key is an enumeration for different auth keys for context.
type Key int

const (
	// KeyTokenClaims is a key to receive auth result from context - claims or error.
	KeyTokenClaims Key = 0
	// KeyToken is a key to receive auth token from context.
	KeyToken Key = 1
)

// SetTokenClaims creates new context with TokenClaims.
func SetTokenClaims(ctx context.Context, claims TokenClaims) context.Context {
	return context.WithValue(ctx, KeyTokenClaims, claims)
}

// GetTokenClaims gets TokenClaims from context.
func GetTokenClaims(ctx context.Context) (TokenClaims, error) {
	value := ctx.Value(KeyTokenClaims)

	if auth, ok := value.(TokenClaims); ok {
		return auth, nil
	}

	if err, ok := value.(error); ok {
		return TokenClaims{}, Error.Wrap(err)
	}

	return TokenClaims{}, Error.New("could not get auth or error from context")
}

// SetToken creates context with auth token.
func SetToken(ctx context.Context, key []byte) context.Context {
	return context.WithValue(ctx, KeyToken, key)
}

// GetToken returns auth token from context is exists.
func GetToken(ctx context.Context) ([]byte, bool) {
	key, ok := ctx.Value(KeyToken).([]byte)
	return key, ok
}
