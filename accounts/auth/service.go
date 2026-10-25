package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/accounts"
	"generator/internal/logger"
	"generator/pkg/auth"
	"generator/pkg/ethsig"
	"generator/server/emails"
)

const (
	// TokenExpirationTime after passing this time token expires.
	TokenExpirationTime = 24 * time.Hour
)

var (
	// ErrUnauthenticated should be returned when user performs unauthenticated action.
	ErrUnauthenticated = errs.Class("user unauthenticated error")

	// Error is a error class for internal auth errors.
	Error = errs.Class("user auth internal error")

	// ErrPermission should be returned when user permission denied.
	ErrPermission = errs.Class("permission denied")

	// ErrAddressAlreadyInUse should be returned when users email address is already in use.
	ErrAddressAlreadyInUse = errs.Class("email address is already in use")

	// ErrAccountExists should be returned when the wallet is registered already.
	ErrAccountExists = errs.Class("account already exists")
)

// Service is handling all user authentication logic.
//
// architecture: Service
type Service struct {
	users        accounts.DB
	usersService *accounts.Service
	signer       auth.JWTSigner
	emailService *emails.Service
	log          logger.Logger
}

// NewService is a constructor for user auth service.
func NewService(users accounts.DB, usersService *accounts.Service, signer auth.JWTSigner, emails *emails.Service,
	log logger.Logger) *Service {
	return &Service{
		users:        users,
		usersService: usersService,
		signer:       signer,
		emailService: emails,
		log:          log,
	}
}

// Authorize validates token from context and returns claims.
func (service *Service) Authorize(ctx context.Context, tokenS string) (auth.TokenClaims, error) {
	token, err := auth.FromBase64URLString(tokenS)
	if err != nil {
		return nil, ErrUnauthenticated.Wrap(err)
	}

	claims, err := service.authenticate(token)
	if err != nil {
		return nil, ErrUnauthenticated.Wrap(err)
	}

	err = service.authorize(ctx, claims)
	if err != nil {
		return nil, ErrUnauthenticated.Wrap(err)
	}

	return claims, nil
}

// authenticate validates token signature and returns authenticated claims.
func (service *Service) authenticate(token auth.Token) (auth.TokenClaims, error) {
	signature := token.Signature

	err := service.signer.SignToken(&token)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	if subtle.ConstantTimeCompare(signature, token.Signature) != 1 {
		return nil, Error.New("incorrect signature")
	}

	claims, err := auth.FromJSON(token.Payload)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return *claims, nil
}

// authorize checks claims and returns authorized User.
func (service *Service) authorize(ctx context.Context, claimsMap auth.TokenClaims) error {
	claims, err := GetStructClaims(claimsMap)
	if err != nil {
		return ErrUnauthenticated.Wrap(err)
	}

	// every token is issued with an expiry time, so a token without it is rejected as well.
	if !claims.ExpiresAt.After(time.Now()) {
		return ErrUnauthenticated.New("token expired")
	}

	_, err = service.users.GetByWalletAddress(ctx, claims.Wallet)
	if err != nil {
		return ErrUnauthenticated.New("authorization failed. no user with wallet address: %s", claims.Wallet)
	}

	return nil
}

// Nonce creates nonce and send to metamask for login.
func (service *Service) Nonce(ctx context.Context, address common.Address) (string, error) {
	user, err := service.users.GetByWalletAddress(ctx, address)
	if err != nil {
		return "", Error.Wrap(err)
	}

	nonce := hexutil.Encode(user.Nonce)

	return nonce, nil
}

// RegisterWithMetamask creates user by credentials.
func (service *Service) RegisterWithMetamask(ctx context.Context, signature []byte) error {
	walletAddress, err := recoverWalletAddress([]byte(accounts.DefaultMessageForRegistration), signature)
	if err != nil {
		return ErrUnauthenticated.Wrap(err)
	}

	_, err = service.users.GetByWalletAddress(ctx, walletAddress)
	switch {
	case err == nil:
		return ErrAccountExists.New("wallet %s", walletAddress.Hex())
	case !accounts.ErrNoAccount.Has(err):
		return Error.Wrap(err)
	}

	nonce := make([]byte, 32)
	_, err = rand.Read(nonce)
	if err != nil {
		return Error.Wrap(err)
	}

	user := accounts.Account{
		ID:        uuid.New(),
		Email:     "",
		Nonce:     nonce,
		LastLogin: time.Time{},
		Status:    accounts.StatusActive,
		CreatedAt: time.Now().UTC(),
		Wallet:    walletAddress,
	}
	err = service.users.Create(ctx, user)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}

// LoginWithMetamask authenticates user by credentials and returns login token.
func (service *Service) LoginWithMetamask(ctx context.Context, nonce string, signature []byte) (string, error) {
	walletAddress, err := recoverWalletAddress([]byte(nonce), signature)
	if err != nil {
		return "", ErrUnauthenticated.Wrap(err)
	}

	user, err := service.users.GetByWalletAddress(ctx, walletAddress)
	if err != nil {
		if accounts.ErrNoAccount.Has(err) {
			return "", ErrUnauthenticated.New("the message is not signed by a registered wallet")
		}
		return "", Error.Wrap(err)
	}

	decodeNonce, err := hexutil.Decode(nonce)
	if err != nil {
		return "", ErrUnauthenticated.Wrap(err)
	}

	// the nonce changes after every login, so a signature can not be used twice.
	if !bytes.Equal(decodeNonce, user.Nonce) {
		return "", ErrUnauthenticated.New("nonce is invalid")
	}

	claims := Claims{
		UserID:    user.ID,
		Wallet:    user.Wallet,
		ExpiresAt: time.Now().Add(TokenExpirationTime),
	}

	goAuthClaims, err := GetMapClaims(claims)
	if err != nil {
		return "", Error.Wrap(err)
	}

	token, err := service.signer.CreateToken(ctx, &goAuthClaims)
	if err != nil {
		return "", Error.Wrap(err)
	}

	newNonce := make([]byte, 32)
	_, err = rand.Read(newNonce)
	if err != nil {
		return "", Error.Wrap(err)
	}

	err = service.users.UpdateNonce(ctx, user.ID, newNonce)
	if err != nil {
		return "", Error.Wrap(err)
	}

	err = service.users.UpdateLastLogin(ctx, user.ID, time.Now().UTC())
	if err != nil {
		service.log.Error("could not update last login", Error.Wrap(err))
	}

	return token, nil
}

// recoverWalletAddress function that verifies the authenticity of the address.
func recoverWalletAddress(message, signature []byte) (common.Address, error) {
	if len(signature) != crypto.SignatureLength {
		return common.Address{}, Error.New("signature must be %d bytes long, got %d", crypto.SignatureLength, len(signature))
	}
	if signature[crypto.RecoveryIDOffset] != 27 && signature[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, Error.New("signature recovery byte is invalid")
	}
	// wallets put 27 or 28 into the recovery byte, SigToPub expects 0 or 1; the caller's slice stays intact.
	sig := append([]byte(nil), signature...)
	sig[crypto.RecoveryIDOffset] -= 27

	pubKey, err := crypto.SigToPub(ethsig.SignHash(message), sig)
	if err != nil {
		return common.Address{}, Error.Wrap(err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return recoveredAddr, nil
}
