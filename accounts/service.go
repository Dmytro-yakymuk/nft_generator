package accounts

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
	"golang.org/x/crypto/bcrypt"
)

// ErrAccounts indicates that there was an error in the service.
var ErrAccounts = errs.Class("accounts service error")

// Service is handling accounts related logic.
//
// architecture: Service
type Service struct {
	accounts DB
}

// NewService is a constructor for accounts service.
func NewService(accounts DB) *Service {
	return &Service{
		accounts: accounts,
	}
}

// Create creates account.
func (service *Service) Create(ctx context.Context, email string, password []byte) error {
	passwordHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return ErrAccounts.Wrap(err)
	}

	account := Account{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		LastLogin:    time.Time{},
		Status:       StatusCreated,
		CreatedAt:    time.Now().UTC(),
	}

	return ErrAccounts.Wrap(service.accounts.Create(ctx, account))
}

// List returns all accounts.
func (service *Service) List(ctx context.Context) ([]Account, error) {
	allAccounts, err := service.accounts.List(ctx)
	return allAccounts, ErrAccounts.Wrap(err)
}

// Get returns account by id.
func (service *Service) Get(ctx context.Context, id uuid.UUID) (Account, error) {
	account, err := service.accounts.Get(ctx, id)
	return account, ErrAccounts.Wrap(err)
}

// Update updates account in the database.
func (service *Service) Update(ctx context.Context, settings AccountSettings) error {
	return ErrAccounts.Wrap(service.accounts.Update(ctx, settings))
}

// UpdateStatus updates status of certain account.
func (service *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	return ErrAccounts.Wrap(service.accounts.UpdateStatus(ctx, id, status))
}

// Delete deletes account.
func (service *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return ErrAccounts.Wrap(service.accounts.Delete(ctx, id))
}
