package auth

import (
	"bytes"
	"encoding/json"

	"github.com/zeebo/errs"
)

// TokenClaimsError is an error class for Claims errors.
var TokenClaimsError = errs.Class("admin auth claims error")

// TokenClaims represents data signed by server and used for authentication.
type TokenClaims map[string]string

// JSON returns json representation of Claims.
func (c *TokenClaims) JSON() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	err := json.NewEncoder(buffer).Encode(c)
	return buffer.Bytes(), TokenClaimsError.Wrap(err)
}

// FromJSON returns TokenClaims instance, parsed from JSON.
func FromJSON(data []byte) (*TokenClaims, error) {
	claims := new(TokenClaims)

	err := json.NewDecoder(bytes.NewReader(data)).Decode(claims)
	if err != nil {
		return nil, TokenClaimsError.Wrap(err)
	}

	return claims, nil
}
