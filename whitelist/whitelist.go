package whitelist

import (
	"bytes"
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/pkg/merkletree"
	"generator/pkg/pagination"
)

// ErrNoItem indicated that item does not exist.
var ErrNoItem = errs.Class("item does not exist")

// DB exposes access to whitelist db.
//
// architecture: DB
type DB interface {
	// Add adds wallet address to whitelist in the database.
	Add(ctx context.Context, projectID uuid.UUID, item Item) error
	// AddMany adds the items to whitelist in one transaction: either all of them or none.
	AddMany(ctx context.Context, projectID uuid.UUID, items []Item) error
	// Get returns item of whitelist by project id and wallet address from the database.
	Get(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) (Item, error)
	// ListByProjectID returns wallet addresses by project id from the database.
	ListByProjectID(ctx context.Context, projectID uuid.UUID) (_ []Item, err error)
	// ListByProjectIDWithPagination returns names and wallet addresses by project id with pagination from the database.
	ListByProjectIDWithPagination(ctx context.Context, projectID uuid.UUID, cursor pagination.Cursor) (_ Page, err error)
	// ListByProjectIDAndItemName returns items by project id and item name from database.
	ListByProjectIDAndItemName(ctx context.Context, projectID uuid.UUID, name string, cursor pagination.Cursor) (_ Page, err error)
	// TotalCount counts items of whitelist in the table.
	TotalCount(ctx context.Context, projectID uuid.UUID) (int, error)
	// Delete deletes wallet address from whitelist in the database.
	Delete(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) error
}

// Content entity describes wallet address required to create merkre tree.
type Content struct {
	WalletAddress common.Address `json:"walletAddress"`
}

// Item entity describes values required to create item of whitelist.
type Item struct {
	Name          string         `json:"name"`
	WalletAddress common.Address `json:"walletAddress"`
}

// Page holds whitelist page entity which is used to show listed page of items.
type Page struct {
	Items []Item          `json:"items"`
	Page  pagination.Page `json:"page"`
}

// Filter defines the list of possible filters.
type Filter string

const (
	// FilterName indicates filtering by item name.
	FilterName Filter = "name"
)

// CalculateHash hashes the wallet address.
func (content Content) CalculateHash() ([]byte, error) {
	return crypto.Keccak256(content.WalletAddress.Bytes()), nil
}

// Equals compares two wallet addresses.
func (content Content) Equals(otherContent merkletree.Content) (bool, error) {
	if res := bytes.Compare(content.WalletAddress.Bytes(), otherContent.(Content).WalletAddress.Bytes()); res == 0 {
		return true, nil
	}
	return false, nil
}

// Config defines configuration for whitelist.
type Config struct {
	Pagination pagination.Cursor `json:"pagination"`
}
