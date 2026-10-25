package whitelist_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator"
	"generator/accounts"
	"generator/collections"
	"generator/collections/layers"
	"generator/collections/layers/filestore"
	"generator/internal/filestore/localfilestore"
	"generator/internal/logger/zaplog"
	"generator/pkg/pagination"
	"generator/storage/testdb"
	"generator/whitelist"
)

func TestWhiteListDB(t *testing.T) {
	user1 := accounts.Account{
		ID:           uuid.New(),
		Email:        "testUser1@gmail.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusCreated,
		CreatedAt:    time.Now().UTC(),
	}

	project1 := collections.Collection{
		ID:     uuid.New(),
		Name:   "project1",
		UserID: user1.ID,
	}

	item1 := whitelist.Item{
		Name:          "test item one",
		WalletAddress: common.HexToAddress("56f088767d91badc379155290c4205c7b917a36e"),
	}

	item2 := whitelist.Item{
		Name:          "item2",
		WalletAddress: common.HexToAddress("56f088767d91badc379155290c4205c7b917a36d"),
	}

	cursor := pagination.Cursor{
		Limit: 2,
		Page:  1,
	}

	filterItemName := "item"

	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		repositoryUsers := db.Accounts()
		repositoryProjects := db.Collections()
		repositoryWhiteList := db.WhiteList()

		t.Run("Create", func(t *testing.T) {
			err := repositoryUsers.Create(ctx, user1)
			require.NoError(t, err)
			err = repositoryProjects.Create(ctx, project1)
			require.NoError(t, err)
			err = repositoryWhiteList.Add(ctx, project1.ID, item1)
			require.NoError(t, err)
			err = repositoryWhiteList.Add(ctx, project1.ID, item2)
			require.NoError(t, err)
		})

		t.Run("Get", func(t *testing.T) {
			item, err := repositoryWhiteList.Get(ctx, project1.ID, item1.WalletAddress)
			require.NoError(t, err)

			compareItemOfWhitelist(t, item1, item)
		})

		t.Run("Negative Get", func(t *testing.T) {
			_, err := repositoryWhiteList.Get(ctx, project1.ID, common.Address{})
			require.Error(t, err)

			assert.True(t, whitelist.ErrNoItem.Has(err))
		})

		t.Run("ListByProjectID", func(t *testing.T) {
			walletAddresses, err := repositoryWhiteList.ListByProjectID(ctx, project1.ID)
			require.NoError(t, err)
			require.Equal(t, len(walletAddresses), 2)
			// ordered by address: the order defines the Merkle root.
			compareItemOfWhitelist(t, walletAddresses[0], item2)
			compareItemOfWhitelist(t, walletAddresses[1], item1)
		})

		t.Run("ListByProjectIDWithPagination", func(t *testing.T) {
			itemsPage, err := repositoryWhiteList.ListByProjectIDWithPagination(ctx, project1.ID, cursor)
			require.NoError(t, err)
			require.Equal(t, len(itemsPage.Items), 2)

			// pages are ordered by name.
			compareItemOfWhitelist(t, itemsPage.Items[0], item2)
			compareItemOfWhitelist(t, itemsPage.Items[1], item1)

			require.Equal(t, itemsPage.Page.Offset, 0)
			require.Equal(t, itemsPage.Page.Limit, 2)
			require.Equal(t, itemsPage.Page.CurrentPage, 1)
			require.Equal(t, itemsPage.Page.PageCount, 1)
			require.Equal(t, itemsPage.Page.TotalCount, 2)
		})

		t.Run("TotalCount", func(t *testing.T) {
			totalCount, err := repositoryWhiteList.TotalCount(ctx, project1.ID)
			require.NoError(t, err)
			require.Equal(t, totalCount, 2)
		})

		t.Run("ListByProjectIDAndItemName", func(t *testing.T) {
			itemsPage, err := repositoryWhiteList.ListByProjectIDAndItemName(ctx, project1.ID, filterItemName, cursor)
			assert.NoError(t, err)
			assert.Equal(t, len(itemsPage.Items), 1)
			compareItemOfWhitelist(t, item1, itemsPage.Items[0])
		})

		t.Run("Delete", func(t *testing.T) {
			err := repositoryWhiteList.Delete(ctx, project1.ID, item1.WalletAddress)
			require.NoError(t, err)
		})

		t.Run("AddMany adds all items or none", func(t *testing.T) {
			item3 := whitelist.Item{Name: "item3", WalletAddress: common.HexToAddress("0x0000000000000000000000000000000000000003")}
			item4 := whitelist.Item{Name: "item4", WalletAddress: common.HexToAddress("0x0000000000000000000000000000000000000004")}

			// item2 is already in the whitelist: the whole batch is rolled back, item3 is not added either.
			err := repositoryWhiteList.AddMany(ctx, project1.ID, []whitelist.Item{item3, item2})
			require.Error(t, err)
			_, err = repositoryWhiteList.Get(ctx, project1.ID, item3.WalletAddress)
			assert.True(t, whitelist.ErrNoItem.Has(err))

			err = repositoryWhiteList.AddMany(ctx, project1.ID, []whitelist.Item{item3, item4})
			require.NoError(t, err)
			totalCount, err := repositoryWhiteList.TotalCount(ctx, project1.ID)
			require.NoError(t, err)
			require.Equal(t, 3, totalCount)
		})
	})
}

func TestWhiteListService(t *testing.T) {
	user1 := accounts.Account{
		ID:           uuid.New(),
		Email:        "testUser1@gmail.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusCreated,
		CreatedAt:    time.Now().UTC(),
	}

	project1 := collections.Collection{
		ID:       uuid.New(),
		Name:     "project1",
		UserID:   user1.ID,
		NFTCount: 10,
	}

	item1 := whitelist.Item{
		Name:          "test item one",
		WalletAddress: common.HexToAddress("0x5EE3e1B1484d415B194974ea2Be6A632D02Ce15D"),
	}
	item2 := whitelist.Item{
		Name:          "item2",
		WalletAddress: common.HexToAddress("0x8894bC144F9C6EDd3ac938a38A437E4af4868934"),
	}
	item3 := whitelist.Item{
		Name:          "item3",
		WalletAddress: common.HexToAddress("0xfaBcb77c3B2429076a59BA1Ea7f0ad0Df3Aa736C"),
	}
	testItem1 := whitelist.Item{
		Name:          "wallet addr1",
		WalletAddress: common.HexToAddress("0x56f088767D91badc379155290c4205c7b917a36E"),
	}
	testItem2 := whitelist.Item{
		Name:          "wallet addr2",
		WalletAddress: common.HexToAddress("0x56f088767D91badc379155290c4205c7b917a361"),
	}

	cursor := pagination.Cursor{
		Limit: 2,
		Page:  1,
	}

	filterItemName := "item"

	merkleRootExpected := "27e35a316792d608873f05f56059e5195f402f90e7fd6e2af152549ac686ed84"
	merkleProofExpected := []string{
		"114b216bd8e8b301b9e4ad57fa99dae1944a4071981565e76cc56a1a7846d08e",
		"6b6b244df3a1069765b46c12cf07b6fcd277987d36de896b64c8073becec9906",
	}

	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		repositoryUsers := db.Accounts()
		repositoryProjects := db.Collections()
		repositoryWhiteList := db.WhiteList()
		repositoryLayers := db.Layers()

		fileStoreLayers := filestore.NewStore(filestore.Config{})
		fileStoreProjectImages := localfilestore.NewFileStore(localfilestore.Config{})
		serviceLayers := layers.NewService(zaplog.NewLog(), layers.Config{}, repositoryLayers, nil, fileStoreLayers, fileStoreProjectImages)
		serviceCollections := collections.NewService(collections.Config{}, repositoryProjects, serviceLayers, fileStoreProjectImages, nil, nil)
		serviceWhiteList := whitelist.NewService(whitelist.Config{}, repositoryWhiteList, serviceCollections)

		t.Run("Create", func(t *testing.T) {
			err := repositoryUsers.Create(ctx, user1)
			require.NoError(t, err)
			err = repositoryProjects.Create(ctx, project1)
			require.NoError(t, err)

			err = serviceWhiteList.Add(ctx, project1.ID, item1)
			require.NoError(t, err)
			err = serviceWhiteList.Add(ctx, project1.ID, item2)
			require.NoError(t, err)
			err = serviceWhiteList.Add(ctx, project1.ID, item3)
			require.NoError(t, err)
		})

		t.Run("ListByProjectIDWithPagination", func(t *testing.T) {
			itemsPage, err := serviceWhiteList.ListByProjectIDWithPagination(ctx, project1.ID, cursor)
			require.NoError(t, err)
			require.Equal(t, len(itemsPage.Items), 2)

			// pages are ordered by name: "item2", "item3", "test item one".
			compareItemOfWhitelist(t, itemsPage.Items[0], item2)
			compareItemOfWhitelist(t, itemsPage.Items[1], item3)

			require.Equal(t, itemsPage.Page.Offset, 0)
			require.Equal(t, itemsPage.Page.Limit, 2)
			require.Equal(t, itemsPage.Page.CurrentPage, 1)
			require.Equal(t, itemsPage.Page.PageCount, 2)
			require.Equal(t, itemsPage.Page.TotalCount, 3)
		})

		t.Run("ListByProjectIDAndItemName", func(t *testing.T) {
			itemsPage, err := serviceWhiteList.ListByProjectIDAndItemName(ctx, project1.ID, filterItemName, cursor)
			assert.NoError(t, err)
			assert.Equal(t, len(itemsPage.Items), 1)
			compareItemOfWhitelist(t, item1, itemsPage.Items[0])
		})

		t.Run("Get", func(t *testing.T) {
			item, err := serviceWhiteList.Get(ctx, project1.ID, item1.WalletAddress)
			require.NoError(t, err)

			compareItemOfWhitelist(t, item1, item)
		})

		t.Run("Negative Get", func(t *testing.T) {
			_, err := serviceWhiteList.Get(ctx, project1.ID, common.Address{})
			require.Error(t, err)

			assert.True(t, whitelist.ErrNoItem.Has(err))
		})

		t.Run("GetMerkleRoot", func(t *testing.T) {
			merkleRoot, err := serviceWhiteList.GetMerkleRoot(ctx, project1.ID)
			require.NoError(t, err)

			mr, err := hex.DecodeString(merkleRootExpected)
			require.NoError(t, err)

			require.Equal(t, mr, merkleRoot)
		})

		t.Run("GetMerkleProof", func(t *testing.T) {
			var merkleProofBytes [][]byte

			merkleProof, err := serviceWhiteList.GetMerkleProof(ctx, project1.ID, item1.WalletAddress)
			require.NoError(t, err)

			for _, v := range merkleProofExpected {
				mp, err := hex.DecodeString(v)
				require.NoError(t, err)

				merkleProofBytes = append(merkleProofBytes, mp)
			}

			require.Equal(t, merkleProofBytes, merkleProof)
		})

		t.Run("AddFromExcelFile", func(t *testing.T) {
			path := "./testdata/success.xlsx"
			fileReader, err := os.Open(path)
			require.NoError(t, err)

			err = serviceWhiteList.AddFromExcelFile(ctx, project1.ID, user1.ID, fileReader)
			require.NoError(t, err)

			for _, expected := range []whitelist.Item{testItem1, testItem2} {
				item, err := serviceWhiteList.Get(ctx, project1.ID, expected.WalletAddress)
				require.NoError(t, err)
				compareItemOfWhitelist(t, expected, item)
			}
		})

		t.Run("Negative AddFromExcelFile. Invalid amount of columns in row", func(t *testing.T) {
			path := "./testdata/invalid_amount_of_columns_in_row.xlsx"
			fileReader, err := os.Open(path)
			require.NoError(t, err)

			err = serviceWhiteList.AddFromExcelFile(ctx, project1.ID, user1.ID, fileReader)
			require.Error(t, err)

			assert.True(t, whitelist.ErrInvalidWhitelistFile.Has(err))
		})

		t.Run("Negative AddFromExcelFile. Invalid wallet address", func(t *testing.T) {
			path := "./testdata/invalid_wallet_address.xlsx"
			fileReader, err := os.Open(path)
			require.NoError(t, err)

			err = serviceWhiteList.AddFromExcelFile(ctx, project1.ID, user1.ID, fileReader)
			require.Error(t, err)

			assert.True(t, whitelist.ErrInvalidWhitelistFile.Has(err))
		})

		t.Run("Negative AddFromExcelFile. Duplication of items in database", func(t *testing.T) {
			path := "./testdata/db_duplication.xlsx"
			fileReader, err := os.Open(path)
			require.NoError(t, err)

			err = serviceWhiteList.AddFromExcelFile(ctx, project1.ID, user1.ID, fileReader)
			require.Error(t, err)

			assert.True(t, whitelist.ErrInvalidWhitelistFile.Has(err))
		})

		t.Run("Negative AddFromExcelFile. Duplication of items in file", func(t *testing.T) {
			path := "./testdata/file_duplication.xlsx"
			fileReader, err := os.Open(path)
			require.NoError(t, err)

			err = serviceWhiteList.AddFromExcelFile(ctx, project1.ID, user1.ID, fileReader)
			require.Error(t, err)

			assert.True(t, whitelist.ErrInvalidWhitelistFile.Has(err))
		})
	})
}

func compareItemOfWhitelist(t *testing.T, item, itemToCompare whitelist.Item) {
	assert.Equal(t, item.Name, itemToCompare.Name)
	assert.Equal(t, item.WalletAddress, itemToCompare.WalletAddress)
}

// verifyProof checks the proof the way OpenZeppelin MerkleProof.verify does in the contract:
// pairs are hashed in ascending order.
func verifyProof(proof [][]byte, root []byte, address common.Address) bool {
	computed := crypto.Keccak256(address.Bytes())
	for _, sibling := range proof {
		if bytes.Compare(computed, sibling) <= 0 {
			computed = crypto.Keccak256(computed, sibling)
		} else {
			computed = crypto.Keccak256(sibling, computed)
		}
	}
	return bytes.Equal(computed, root)
}

func TestMerkleRootAndProofs(t *testing.T) {
	user := accounts.Account{
		ID:           uuid.New(),
		Email:        "merkle@example.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusCreated,
		CreatedAt:    time.Now().UTC(),
	}
	project := collections.Collection{ID: uuid.New(), Name: "merkle", UserID: user.ID, NFTCount: 10}
	other := collections.Collection{ID: uuid.New(), Name: "other", UserID: user.ID, NFTCount: 10}

	// inserted out of address order on purpose.
	addresses := []common.Address{
		common.HexToAddress("0xfaBcb77c3B2429076a59BA1Ea7f0ad0Df3Aa736C"),
		common.HexToAddress("0x5EE3e1B1484d415B194974ea2Be6A632D02Ce15D"),
		common.HexToAddress("0x8894bC144F9C6EDd3ac938a38A437E4af4868934"),
	}

	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		fileStoreProjectImages := localfilestore.NewFileStore(localfilestore.Config{})
		serviceLayers := layers.NewService(zaplog.NewLog(), layers.Config{}, db.Layers(), nil, filestore.NewStore(filestore.Config{}), fileStoreProjectImages)
		serviceCollections := collections.NewService(collections.Config{}, db.Collections(), serviceLayers, fileStoreProjectImages, nil, nil)
		service := whitelist.NewService(whitelist.Config{}, db.WhiteList(), serviceCollections)

		require.NoError(t, db.Accounts().Create(ctx, user))
		require.NoError(t, db.Collections().Create(ctx, project))
		require.NoError(t, db.Collections().Create(ctx, other))

		t.Run("empty whitelist has the zero root", func(t *testing.T) {
			root, err := service.GetMerkleRoot(ctx, project.ID)
			require.NoError(t, err)
			assert.Equal(t, make([]byte, 32), root)

			_, err = service.GetMerkleProof(ctx, project.ID, addresses[0])
			assert.True(t, whitelist.ErrNoItem.Has(err))
		})

		t.Run("one address: the root is its leaf, the proof is empty", func(t *testing.T) {
			require.NoError(t, service.Add(ctx, project.ID, whitelist.Item{Name: "one", WalletAddress: addresses[0]}))

			root, err := service.GetMerkleRoot(ctx, project.ID)
			require.NoError(t, err)
			assert.Equal(t, crypto.Keccak256(addresses[0].Bytes()), root)

			proof, err := service.GetMerkleProof(ctx, project.ID, addresses[0])
			require.NoError(t, err)
			assert.Empty(t, proof)
			assert.True(t, verifyProof(proof, root, addresses[0]))
		})

		t.Run("the root does not depend on the order of rows", func(t *testing.T) {
			for i, address := range addresses[1:] {
				require.NoError(t, service.Add(ctx, project.ID, whitelist.Item{Name: "item " + strconv.Itoa(i), WalletAddress: address}))
			}
			root, err := service.GetMerkleRoot(ctx, project.ID)
			require.NoError(t, err)
			for _, address := range addresses {
				proof, err := service.GetMerkleProof(ctx, project.ID, address)
				require.NoError(t, err)
				assert.True(t, verifyProof(proof, root, address), address.Hex())
			}

			// deleting and adding the address again changes how rows are stored, but not the root.
			require.NoError(t, service.Delete(ctx, project.ID, addresses[0]))
			require.NoError(t, service.Add(ctx, project.ID, whitelist.Item{Name: "one", WalletAddress: addresses[0]}))
			again, err := service.GetMerkleRoot(ctx, project.ID)
			require.NoError(t, err)
			assert.Equal(t, root, again)
		})

		t.Run("search does not return addresses of other collections", func(t *testing.T) {
			require.NoError(t, service.Add(ctx, other.ID, whitelist.Item{Name: "item of other", WalletAddress: addresses[0]}))

			page, err := service.ListByProjectIDAndItemName(ctx, project.ID, "item", pagination.Cursor{Limit: 10, Page: 1})
			require.NoError(t, err)
			assert.Len(t, page.Items, 2)
			assert.Equal(t, 2, page.Page.TotalCount)
			for _, item := range page.Items {
				assert.NotEqual(t, "item of other", item.Name)
			}
		})
	})
}
