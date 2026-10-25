package auth

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"generator/pkg/auth"
)

// Claims represents data signed by server and used for authentication.
type Claims struct {
	UserID    uuid.UUID      `json:"userId"`
	Wallet    common.Address `json:"wallet"`
	ExpiresAt time.Time      `json:"expiresAt"`
}

// GetStructClaims returns claims in struct from auth.
func GetStructClaims(claimsMap auth.TokenClaims) (*Claims, error) {
	var claims Claims

	jsonClaims, err := json.Marshal(claimsMap)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonClaims, &claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

// GetMapClaims returns claims in map for goauth from Claims struct.
func GetMapClaims(claims Claims) (auth.TokenClaims, error) {
	var goAuthClaims auth.TokenClaims = make(map[string]string)

	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonClaims, &goAuthClaims)
	if err != nil {
		return nil, err
	}

	return goAuthClaims, nil
}
