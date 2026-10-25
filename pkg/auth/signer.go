package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/zeebo/errs"
)

// JWTSignerError is an error class for JWTSigner errors.
var JWTSignerError = errs.Class("auth token signer error")

// JWTSigner creates signature for provided auth Token. Its hmac256 based JWTSigner.
type JWTSigner struct {
	Secret []byte
}

// SignToken signs token with some secret key.
func (a *JWTSigner) SignToken(token *Token) error {
	mac := hmac.New(sha256.New, a.Secret)

	encoded := base64.URLEncoding.EncodeToString(token.Payload)

	_, err := mac.Write([]byte(encoded))
	if err != nil {
		return JWTSignerError.Wrap(err)
	}

	token.Signature = mac.Sum(nil)

	return nil
}

// CreateToken creates string representation.
func (a *JWTSigner) CreateToken(ctx context.Context, claims *TokenClaims) (string, error) {
	json, err := claims.JSON()
	if err != nil {
		return "", JWTSignerError.Wrap(err)
	}

	token := Token{Payload: json}
	err = a.SignToken(&token)
	if err != nil {
		return "", JWTSignerError.Wrap(err)
	}

	return token.String(), nil
}
