package accounts_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"generator"
	"generator/accounts"
	"generator/accounts/auth"
	"generator/internal/logger"
	pkgauth "generator/pkg/auth"
	"generator/server/emails"
	"generator/storage/testdb"
)

func TestAccountsDB(t *testing.T) {
	address1 := "0xb2cdC7EB2F9d2E629ee97BB91700622A42e688b7"
	address2 := "0xb2cdC7EB2F9d2E629ee97BB91700622A42e688b8"
	user1 := accounts.Account{
		ID:           uuid.New(),
		Name:         "user1",
		Wallet:       common.HexToAddress(address1),
		Email:        "testUser1@gmail.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusCreated,
		IsSubscribed: true,
		CreatedAt:    time.Now().UTC(),
	}

	user2 := accounts.Account{
		ID:           uuid.New(),
		Name:         "user2",
		Wallet:       common.HexToAddress(address2),
		Email:        "testUser1@gmail.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusActive,
		IsSubscribed: false,
		CreatedAt:    time.Now().UTC(),
	}

	userSetting1 := accounts.AccountSettings{
		ID:           user1.ID,
		Name:         "user3",
		Email:        "testUser1@gmail.com3",
		Status:       accounts.StatusCreated,
		IsSubscribed: true,
	}

	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		repositoryAccounts := db.Accounts()

		t.Run("get sql no rows", func(t *testing.T) {
			_, err := repositoryAccounts.Get(ctx, uuid.New())
			require.Error(t, err)
			assert.Equal(t, true, accounts.ErrNoAccount.Has(err))
		})

		t.Run("create", func(t *testing.T) {
			err := repositoryAccounts.Create(ctx, user1)
			require.NoError(t, err)
		})

		t.Run("get", func(t *testing.T) {
			userFromDB, err := repositoryAccounts.Get(ctx, user1.ID)
			require.NoError(t, err)
			compareUsers(t, user1, userFromDB)
		})

		t.Run("get by email sql no rows", func(t *testing.T) {
			_, err := repositoryAccounts.GetByEmail(ctx, "fakemail@gmail.com")
			require.Error(t, err)
			assert.Equal(t, true, accounts.ErrNoAccount.Has(err))
		})

		t.Run("get by email", func(t *testing.T) {
			userFromDB, err := repositoryAccounts.GetByEmail(ctx, user1.Email)
			require.NoError(t, err)
			compareUsers(t, user1, userFromDB)
		})

		t.Run("get by wallet address", func(t *testing.T) {
			userFromDB, err := repositoryAccounts.GetByWalletAddress(ctx, user1.Wallet)
			require.NoError(t, err)
			compareUsers(t, user1, userFromDB)
		})

		t.Run("list", func(t *testing.T) {
			err := repositoryAccounts.Create(ctx, user2)
			require.NoError(t, err)

			allUsers, err := repositoryAccounts.List(ctx)
			require.NoError(t, err)
			require.Equal(t, len(allUsers), 2)
			compareUsers(t, user1, allUsers[0])
			compareUsers(t, user2, allUsers[1])
		})

		t.Run("Update", func(t *testing.T) {
			user1.Name = userSetting1.Name
			user1.Email = userSetting1.Email
			user1.Status = userSetting1.Status
			user1.IsSubscribed = userSetting1.IsSubscribed
			err := repositoryAccounts.Update(ctx, userSetting1)
			require.NoError(t, err)

			userFromDB, err := repositoryAccounts.GetByWalletAddress(ctx, user1.Wallet)
			require.NoError(t, err)
			compareUsers(t, user1, userFromDB)
		})

		t.Run("update status sql no rows", func(t *testing.T) {
			err := repositoryAccounts.UpdateStatus(ctx, uuid.New(), accounts.StatusActive)
			require.Error(t, err)
			require.Equal(t, accounts.ErrNoAccount.Has(err), true)
		})

		t.Run("update status", func(t *testing.T) {
			err := repositoryAccounts.UpdateStatus(ctx, user1.ID, accounts.StatusActive)
			require.NoError(t, err)

			userFromDB, err := repositoryAccounts.Get(ctx, user1.ID)
			require.NoError(t, err)
			assert.Equal(t, accounts.StatusActive, userFromDB.Status)
		})

		t.Run("update password hash sql no rows", func(t *testing.T) {
			err := repositoryAccounts.UpdatePasswordHash(ctx, uuid.New(), []byte("newPass123_"))
			require.Error(t, err)
			require.Equal(t, true, accounts.ErrNoAccount.Has(err))
		})

		t.Run("update password hash", func(t *testing.T) {
			newPasswordHash := []byte("newPass123_")
			err := repositoryAccounts.UpdatePasswordHash(ctx, user1.ID, newPasswordHash)
			require.NoError(t, err)

			userFromDB, err := repositoryAccounts.Get(ctx, user1.ID)
			require.NoError(t, err)
			assert.Equal(t, newPasswordHash, userFromDB.PasswordHash)
		})

		t.Run("delete sql no rows", func(t *testing.T) {
			err := repositoryAccounts.Delete(ctx, uuid.New())
			require.Error(t, err)
			require.Equal(t, accounts.ErrNoAccount.Has(err), true)
		})

		t.Run("delete", func(t *testing.T) {
			err := repositoryAccounts.Delete(ctx, user1.ID)
			require.NoError(t, err)
		})
	})
}

func TestAccountService(t *testing.T) {
	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		var log logger.Logger
		accountService := accounts.NewService(db.Accounts())
		authService := auth.NewService(db.Accounts(), accountService, pkgauth.JWTSigner{}, &emails.Service{}, log)
		wallet := common.HexToAddress("0x2346b33f2e379dda22b2563b009382a0fc9aa926")
		signature, err := hexutil.Decode("0x3b8566bdf04b3e11e825196e6d25c436814c43e40a44db89b17450d286b2e531640309222d2a30a9235f9521b9cba34ca6d601d9e04ec0f75466c33fe2c720b71b")
		require.NoError(t, err)

		testUser1 := accounts.Account{
			Email:        "testUser@gmail.com",
			PasswordHash: []byte("my-pass"),
		}

		t.Run("create", func(t *testing.T) {
			err := accountService.Create(ctx, testUser1.Email, testUser1.PasswordHash)
			require.NoError(t, err)
		})

		t.Run("check is password correct", func(t *testing.T) {
			allUsers, err := accountService.List(ctx)
			require.NoError(t, err)
			require.Equal(t, len(allUsers), 1)

			err = bcrypt.CompareHashAndPassword(allUsers[0].PasswordHash, testUser1.PasswordHash)
			require.NoError(t, err)
		})

		t.Run("Register user", func(t *testing.T) {
			err := authService.RegisterWithMetamask(ctx, signature)
			require.NoError(t, err)
		})

		t.Run("get nonce", func(t *testing.T) {
			user, err := db.Accounts().GetByWalletAddress(ctx, wallet)
			require.NoError(t, err)
			userNonce := hexutil.Encode(user.Nonce)

			nonce, err := authService.Nonce(ctx, wallet)
			require.NoError(t, err)
			require.Equal(t, nonce, userNonce)
		})

		t.Run("login", func(t *testing.T) {
			nonce := "0xfd71049a838a6f7983dae92608f85d56f42318d5d481697c714c634be4875315"
			sign := "0x64bfdf19719d55fde70a9f9e3b1087d2429b54ed712142855d9b98d0db7643f840ca438eb9cc6363fbc772a6c11030a5ae30d9b78ab6f19fed742910caa29ab31c"
			userNonce, err := hexutil.Decode(nonce)
			require.NoError(t, err)
			signature, err := hexutil.Decode(sign)
			require.NoError(t, err)

			user, err := db.Accounts().GetByWalletAddress(ctx, wallet)
			require.NoError(t, err)

			err = db.Accounts().UpdateNonce(ctx, user.ID, userNonce)
			require.NoError(t, err)

			_, err = authService.LoginWithMetamask(ctx, nonce, signature)
			require.NoError(t, err)
		})
	})
}

func compareUsers(t *testing.T, user1, user2 accounts.Account) {
	assert.Equal(t, user1.ID, user2.ID)
	assert.Equal(t, user1.Name, user2.Name)
	assert.Equal(t, user1.Email, user2.Email)
	assert.Equal(t, user1.PasswordHash, user2.PasswordHash)
	assert.WithinDuration(t, user1.LastLogin, user2.LastLogin, 1*time.Second)
	assert.Equal(t, user1.Status, user2.Status)
	assert.Equal(t, user1.IsSubscribed, user2.IsSubscribed)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, 1*time.Second)
}
