package accounts

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrNoAccount indicated that account does not exist.
var ErrNoAccount = errs.Class("account does not exist")

// DB exposes access to accounts db.
//
// architecture: DB
type DB interface {
	// Create creates an account and writes to the database.
	Create(ctx context.Context, account Account) error
	// List returns all accounts from the database.
	List(ctx context.Context) ([]Account, error)
	// Get returns account by id from the database.
	Get(ctx context.Context, id uuid.UUID) (Account, error)
	// GetByEmail returns account by email from the database.
	GetByEmail(ctx context.Context, email string) (Account, error)
	// GetByWalletAddress returns account by wallet address from the database.
	GetByWalletAddress(ctx context.Context, walletAddress common.Address) (Account, error)
	// Update updates account in the database.
	Update(ctx context.Context, settings AccountSettings) error
	// UpdateStatus updates account status in the database.
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	// UpdatePasswordHash updates a password hash in the database.
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte) error
	// UpdateLastLogin updates last login time in the database.
	UpdateLastLogin(ctx context.Context, id uuid.UUID, lastLogin time.Time) error
	// UpdateNonce creates a nonce and writes to the database.
	UpdateNonce(ctx context.Context, id uuid.UUID, nonce []byte) error
	// Delete deletes an account in the database.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Status defines the list of possible user statuses.
type Status int

const (
	// StatusCreated indicates that user is created but not confirmed email yet.
	StatusCreated Status = 0
	// StatusActive indicates that user confirmed his email.
	StatusActive Status = 1
)

// DefaultMessageForRegistration use for registration user by metamask.
const DefaultMessageForRegistration = "Register with metamask"

// Account describes account entity. The login nonce and the password hash are not sent to clients.
type Account struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	Wallet       common.Address `json:"wallet"`
	Nonce        []byte         `json:"-"`
	Email        string         `json:"email"`
	PasswordHash []byte         `json:"-"`
	LastLogin    time.Time      `json:"lastLogin"`
	Status       Status         `json:"status"`
	IsSubscribed bool           `json:"isSubscribed"`
	CreatedAt    time.Time      `json:"createdAt"`
}

// AccountSettings entity describes values required to change account's fields.
type AccountSettings struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Status       Status    `json:"status"`
	IsSubscribed bool      `json:"isSubscribed"`
}
