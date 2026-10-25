package whitelist

import (
	"bytes"
	"context"
	"fmt"
	"io"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeebo/errs"
	"golang.org/x/crypto/sha3"

	"generator/collections"
	"generator/pkg/merkletree"
	"generator/pkg/pagination"
)

var (
	// ErrWhiteList indicates that there was an error in the service.
	ErrWhiteList = errs.Class("whitelist service error")

	// ErrInvalidWhitelistFile indicates that there was an error in the whitelist file.
	ErrInvalidWhitelistFile = errs.Class("invalid whitelist file error")
)

// isDuplicate indicates that the elements of the merkel tree are not duplicated.
const isDuplicate bool = false

// isSortChash indicates that every two elements of the Merkel tree are sorted in ascending order.
const isSortChash bool = true

// isHashOne indicates that the last, single element of the Merkel tree is not hashish.
const isHashOne bool = false

// Service is handling whitelist related logic.
//
// architecture: Service
type Service struct {
	config      Config
	whitelist   DB
	collections *collections.Service
}

// NewService is a constructor for whitelist service.
func NewService(config Config, whitelist DB, collections *collections.Service) *Service {
	return &Service{
		config:      config,
		whitelist:   whitelist,
		collections: collections,
	}
}

// Add adds wallet address to whitelist in the database.
func (service *Service) Add(ctx context.Context, projectID uuid.UUID, item Item) error {
	return ErrWhiteList.Wrap(service.whitelist.Add(ctx, projectID, item))
}

// Get returns item of whitelist by project id and wallet address from the database.
func (service *Service) Get(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) (Item, error) {
	item, err := service.whitelist.Get(ctx, projectID, walletAddress)
	return item, ErrWhiteList.Wrap(err)
}

// AddFromExcelFile adds wallet addresses to whitelist from excel file in the database. Every row is checked
// first, and the addresses are added in one transaction, so a bad file or a failed insert adds nothing.
func (service *Service) AddFromExcelFile(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, fileReader io.Reader) error {
	project, err := service.collections.Get(ctx, projectID)
	if err != nil {
		return ErrWhiteList.Wrap(err)
	}
	if project.UserID != userID {
		return collections.ErrCollectionOwnership.New("")
	}

	file, err := excelize.OpenReader(fileReader)
	if err != nil {
		return ErrWhiteList.Wrap(err)
	}

	var items []Item
	seen := make(map[common.Address]bool)
	for _, sheet := range file.GetSheetList() {
		rows, err := file.GetRows(sheet)
		if err != nil {
			return ErrWhiteList.Wrap(err)
		}

		for _, row := range rows {
			if len(row) != 2 {
				return ErrInvalidWhitelistFile.New("row does not have 2 columns")
			}

			if !common.IsHexAddress(row[1]) {
				return ErrInvalidWhitelistFile.New("invalid wallet address")
			}

			walletAddress := common.HexToAddress(row[1])
			item, err := service.whitelist.Get(ctx, projectID, walletAddress)
			if err != nil {
				if !ErrNoItem.Has(err) {
					return ErrWhiteList.Wrap(err)
				}
			}
			if item.WalletAddress != common.HexToAddress("") {
				return ErrInvalidWhitelistFile.New(fmt.Sprintf("duplication of items in database starting with item where wallet address equal %v",
					walletAddress))
			}

			if seen[walletAddress] {
				return ErrInvalidWhitelistFile.New(fmt.Sprintf("duplication of items in file where wallet address equal %v", walletAddress))
			}
			seen[walletAddress] = true

			items = append(items, Item{Name: row[0], WalletAddress: walletAddress})
		}
	}

	return ErrWhiteList.Wrap(service.whitelist.AddMany(ctx, projectID, items))
}

// ListByProjectID returns wallet addresses by project id from the database.
func (service *Service) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]Item, error) {
	items, err := service.whitelist.ListByProjectID(ctx, projectID)
	return items, ErrWhiteList.Wrap(err)
}

// ListByProjectIDWithPagination returns names and wallet addresses by project id with pagination from the database.
func (service *Service) ListByProjectIDWithPagination(ctx context.Context, projectID uuid.UUID, cursor pagination.Cursor) (_ Page, err error) {
	if cursor.Limit <= 0 {
		cursor.Limit = service.config.Pagination.Limit
	}
	if cursor.Page <= 0 {
		cursor.Page = service.config.Pagination.Page
	}
	items, err := service.whitelist.ListByProjectIDWithPagination(ctx, projectID, cursor)
	return items, ErrWhiteList.Wrap(err)
}

// ListByProjectIDAndItemName returns items by project id and item name from database.
func (service *Service) ListByProjectIDAndItemName(ctx context.Context, projectID uuid.UUID, name string, cursor pagination.Cursor) (_ Page, err error) {
	if cursor.Limit <= 0 {
		cursor.Limit = service.config.Pagination.Limit
	}
	if cursor.Page <= 0 {
		cursor.Page = service.config.Pagination.Page
	}
	items, err := service.whitelist.ListByProjectIDAndItemName(ctx, projectID, name, cursor)
	return items, ErrWhiteList.Wrap(err)
}

// GetMerkleRoot returns the Merkle root of the collection whitelist. The addresses are read
// in a fixed order, so the root does not change until the whitelist does; an empty whitelist
// has the zero root, with which nobody can mint through the whitelist.
func (service *Service) GetMerkleRoot(ctx context.Context, projectID uuid.UUID) ([]byte, error) {
	var list []merkletree.Content

	items, err := service.whitelist.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, ErrWhiteList.Wrap(err)
	}
	if len(items) == 0 {
		return make([]byte, 32), nil
	}

	for _, item := range items {
		list = append(list, Content{WalletAddress: item.WalletAddress})
	}

	tree, err := merkletree.NewTreeWithParameters(list, sha3.NewLegacyKeccak256, isDuplicate, isSortChash, isHashOne)
	if err != nil {
		return nil, ErrWhiteList.Wrap(err)
	}

	return tree.MerkleRoot(), nil
}

// GetMerkleProof returns merkle proof for wallet address.
func (service *Service) GetMerkleProof(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) ([][]byte, error) {
	var (
		list                 []merkletree.Content
		isAddressInWhitelist bool
	)

	items, err := service.whitelist.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, ErrWhiteList.Wrap(err)
	}
	if len(items) == 0 {
		return nil, ErrNoItem.New("whitelist is empty")
	}

	for _, item := range items {
		list = append(list, Content{WalletAddress: item.WalletAddress})
		if res := bytes.Compare(item.WalletAddress.Bytes(), walletAddress.Bytes()); res == 0 {
			isAddressInWhitelist = true
		}
	}
	if !isAddressInWhitelist {
		return nil, ErrNoItem.New("address %s is not in the whitelist", walletAddress.Hex())
	}

	tree, err := merkletree.NewTreeWithParameters(list, sha3.NewLegacyKeccak256, isDuplicate, isSortChash, isHashOne)
	if err != nil {
		return nil, ErrWhiteList.Wrap(err)
	}

	merkleProof, _, err := tree.GetProof(Content{WalletAddress: walletAddress})
	return merkleProof, ErrWhiteList.Wrap(err)
}

// Delete deletes wallet address from whitelist in the database.
func (service *Service) Delete(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) error {
	return ErrWhiteList.Wrap(service.whitelist.Delete(ctx, projectID, walletAddress))
}

// GetPublicMerkleProof returns Merkle proof for a buyer on the published collection mint site.
func (service *Service) GetPublicMerkleProof(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) ([][]byte, error) {
	project, err := service.collections.Get(ctx, projectID)
	if err != nil {
		return nil, ErrWhiteList.Wrap(err)
	}
	if !project.IsWebsitePublished {
		return nil, collections.ErrNoCollection.New("website of collection %s is not published", projectID)
	}

	return service.GetMerkleProof(ctx, projectID, walletAddress)
}
