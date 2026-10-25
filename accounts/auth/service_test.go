package auth_test

import (
	"context"
	"crypto/ecdsa"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator"
	"generator/accounts"
	"generator/accounts/auth"
	"generator/internal/logger/zaplog"
	pkgauth "generator/pkg/auth"
	"generator/pkg/ethsig"
	"generator/server/emails"
	"generator/storage/testdb"
)

// signMessage signs the message the way MetaMask does (personal_sign, recovery byte 27 or 28).
func signMessage(t *testing.T, key *ecdsa.PrivateKey, message string) []byte {
	signature, err := crypto.Sign(ethsig.SignHash([]byte(message)), key)
	require.NoError(t, err)
	signature[crypto.RecoveryIDOffset] += 27
	return signature
}

func TestMetamaskAuth(t *testing.T) {
	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		signer := pkgauth.JWTSigner{Secret: []byte("test secret")}
		service := auth.NewService(db.Accounts(), accounts.NewService(db.Accounts()), signer, &emails.Service{}, zaplog.NewLog())

		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		wallet := crypto.PubkeyToAddress(key.PublicKey)

		t.Run("register", func(t *testing.T) {
			require.NoError(t, service.RegisterWithMetamask(ctx, signMessage(t, key, accounts.DefaultMessageForRegistration)))

			err := service.RegisterWithMetamask(ctx, signMessage(t, key, accounts.DefaultMessageForRegistration))
			assert.True(t, auth.ErrAccountExists.Has(err), err)
		})

		nonce, err := service.Nonce(ctx, wallet)
		require.NoError(t, err)
		signature := signMessage(t, key, nonce)

		t.Run("malformed signatures are rejected without panic", func(t *testing.T) {
			wrongRecoveryByte := append([]byte(nil), signature...)
			wrongRecoveryByte[crypto.RecoveryIDOffset] = 29

			for name, signature := range map[string][]byte{
				"empty":               nil,
				"short":               signature[:64],
				"long":                append(append([]byte(nil), signature...), 0),
				"wrong recovery byte": wrongRecoveryByte,
			} {
				_, err := service.LoginWithMetamask(ctx, nonce, signature)
				assert.True(t, auth.ErrUnauthenticated.Has(err), "%s: %v", name, err)

				err = service.RegisterWithMetamask(ctx, signature)
				assert.True(t, auth.ErrUnauthenticated.Has(err), "%s: %v", name, err)
			}
		})

		t.Run("signature of an unregistered wallet", func(t *testing.T) {
			stranger, err := crypto.GenerateKey()
			require.NoError(t, err)

			_, err = service.LoginWithMetamask(ctx, nonce, signMessage(t, stranger, nonce))
			assert.True(t, auth.ErrUnauthenticated.Has(err), err)
		})

		t.Run("signature of another nonce", func(t *testing.T) {
			other := "0x" + strings.Repeat("ab", 32)
			_, err := service.LoginWithMetamask(ctx, other, signMessage(t, key, other))
			assert.True(t, auth.ErrUnauthenticated.Has(err), err)
		})

		var token string
		t.Run("login", func(t *testing.T) {
			token, err = service.LoginWithMetamask(ctx, nonce, signature)
			require.NoError(t, err)

			claims, err := service.Authorize(ctx, token)
			require.NoError(t, err)
			assert.Equal(t, strings.ToLower(wallet.Hex()), strings.ToLower(claims["wallet"]))
		})

		t.Run("the same signature can not be used twice", func(t *testing.T) {
			_, err := service.LoginWithMetamask(ctx, nonce, signature)
			assert.True(t, auth.ErrUnauthenticated.Has(err), err)
		})

		user, err := db.Accounts().GetByWalletAddress(ctx, wallet)
		require.NoError(t, err)

		strangerKey, err := crypto.GenerateKey()
		require.NoError(t, err)

		// tokenWith signs claims the way LoginWithMetamask does.
		tokenWith := func(signer pkgauth.JWTSigner, claims auth.Claims) string {
			mapClaims, err := auth.GetMapClaims(claims)
			require.NoError(t, err)
			token, err := signer.CreateToken(ctx, &mapClaims)
			require.NoError(t, err)
			return token
		}
		valid := auth.Claims{UserID: user.ID, Wallet: wallet, ExpiresAt: time.Now().Add(time.Hour)}
		expired, withoutExpiry, unknownWallet := valid, valid, valid
		expired.ExpiresAt = time.Now().Add(-time.Minute)
		withoutExpiry.ExpiresAt = time.Time{}
		unknownWallet.UserID, unknownWallet.Wallet = uuid.New(), crypto.PubkeyToAddress(strangerKey.PublicKey)

		_, err = service.Authorize(ctx, tokenWith(signer, valid))
		require.NoError(t, err)
		for name, token := range map[string]string{
			"garbage":               "not a token",
			"signed by another key": tokenWith(pkgauth.JWTSigner{Secret: []byte("another secret")}, valid),
			"expired":               tokenWith(signer, expired),
			"without expiry":        tokenWith(signer, withoutExpiry),
			"unknown wallet":        tokenWith(signer, unknownWallet),
		} {
			token := token
			t.Run("authorize: "+name, func(t *testing.T) {
				_, err := service.Authorize(ctx, token)
				assert.True(t, auth.ErrUnauthenticated.Has(err), err)
			})
		}
	})
}
