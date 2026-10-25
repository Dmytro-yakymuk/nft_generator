package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/accounts"
	"generator/internal/mail"
)

// ensures that accountsDB implements accounts.DB.
var _ accounts.DB = (*accountsDB)(nil)

// ErrAccountsDB indicates that there was an error in the database.
var ErrAccountsDB = errs.Class("accounts repository error")

// accountsDB provides access to accounts db.
//
// architecture: Database
type accountsDB struct {
	conn *sql.DB
}

// Create creates a user and writes to the database.
func (accountsDB *accountsDB) Create(ctx context.Context, user accounts.Account) error {
	emailNormalized := mail.Normalize(user.Email)
	wallet := strings.ToLower(user.Wallet.String())

	query := `INSERT INTO accounts(id, name, wallet_address, nonce, email, email_normalized, password_hash, last_login, status, is_subscribed, created_at)
 	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := accountsDB.conn.ExecContext(ctx, query, user.ID, user.Name, wallet, user.Nonce, user.Email, emailNormalized, user.PasswordHash,
		user.LastLogin, user.Status, user.IsSubscribed, user.CreatedAt)

	return ErrAccountsDB.Wrap(err)
}

// Get returns user by id from the data base.
func (accountsDB *accountsDB) Get(ctx context.Context, id uuid.UUID) (accounts.Account, error) {
	var user accounts.Account
	var address string
	query := `SELECT id, name, wallet_address, nonce, email, password_hash, last_login, status, is_subscribed, created_at 
 	          FROM accounts
 	          WHERE id = $1`

	err := accountsDB.conn.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &address, &user.Nonce, &user.Email,
		&user.PasswordHash, &user.LastLogin, &user.Status, &user.IsSubscribed, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return user, accounts.ErrNoAccount.Wrap(err)
	}

	user.Wallet = common.HexToAddress(address)

	return user, ErrAccountsDB.Wrap(err)
}

// GetByEmail returns user by email from the data base.
func (accountsDB *accountsDB) GetByEmail(ctx context.Context, email string) (accounts.Account, error) {
	var user accounts.Account
	var address string
	emailNormalized := mail.Normalize(email)
	query := `SELECT id, name, wallet_address, nonce, email, password_hash, last_login, status, is_subscribed, created_at 
 	          FROM accounts
 	          WHERE email_normalized = $1`

	err := accountsDB.conn.QueryRowContext(ctx, query, emailNormalized).Scan(&user.ID, &user.Name, &address, &user.Nonce,
		&user.Email, &user.PasswordHash, &user.LastLogin, &user.Status, &user.IsSubscribed, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return user, accounts.ErrNoAccount.Wrap(err)
	}

	user.Wallet = common.HexToAddress(address)

	return user, ErrAccountsDB.Wrap(err)
}

// GetByWalletAddress returns user by wallet address from the database.
func (accountsDB *accountsDB) GetByWalletAddress(ctx context.Context, walletAddress common.Address) (accounts.Account, error) {
	var (
		user    accounts.Account
		address string
	)

	query := "SELECT id, name, wallet_address, nonce, email, password_hash, last_login, status, is_subscribed, created_at FROM accounts WHERE wallet_address = $1"
	wallet := strings.ToLower(walletAddress.String())

	row := accountsDB.conn.QueryRowContext(ctx, query, wallet)

	err := row.Scan(&user.ID, &user.Name, &address, &user.Nonce, &user.Email, &user.PasswordHash, &user.LastLogin, &user.Status,
		&user.IsSubscribed, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, accounts.ErrNoAccount.New("user does not exist")
		}
		return user, ErrAccountsDB.Wrap(err)
	}

	user.Wallet = common.HexToAddress(address)

	return user, nil
}

// List returns all accounts from the data base.
func (accountsDB *accountsDB) List(ctx context.Context) ([]accounts.Account, error) {
	var allUsers []accounts.Account
	query := `SELECT id, name, wallet_address, nonce, email, password_hash, last_login, status, is_subscribed, created_at
 	          FROM accounts`

	rows, err := accountsDB.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrAccountsDB.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var user accounts.Account
		var address string

		err := rows.Scan(&user.ID, &user.Name, &address, &user.Nonce, &user.Email, &user.PasswordHash, &user.LastLogin, &user.Status,
			&user.IsSubscribed, &user.CreatedAt)
		if err != nil {
			return nil, ErrAccountsDB.Wrap(err)
		}

		user.Wallet = common.HexToAddress(address)

		allUsers = append(allUsers, user)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrAccountsDB.Wrap(err)
	}

	return allUsers, nil
}

// Update updates user in the database.
func (accountsDB *accountsDB) Update(ctx context.Context, userSetting accounts.AccountSettings) error {
	emailNormalized := mail.Normalize(userSetting.Email)
	result, err := accountsDB.conn.ExecContext(ctx, "UPDATE accounts SET name=$1, email=$2, email_normalized=$3, status=$4, is_subscribed=$5 WHERE id=$6",
		userSetting.Name, userSetting.Email, emailNormalized, userSetting.Status, userSetting.IsSubscribed, userSetting.ID)
	if err != nil {
		return ErrAccountsDB.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return accounts.ErrNoAccount.New("user does not exist")
	}

	return ErrAccountsDB.Wrap(err)
}

// UpdatePasswordHash updates a password hash in the database.
func (accountsDB *accountsDB) UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte) error {
	result, err := accountsDB.conn.ExecContext(ctx, "UPDATE accounts SET password_hash=$1 WHERE id=$2", passwordHash, id)
	if err != nil {
		return ErrAccountsDB.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return accounts.ErrNoAccount.New("user does not exist")
	}

	return ErrAccountsDB.Wrap(err)
}

// UpdateLastLogin updates last login time.
func (accountsDB *accountsDB) UpdateLastLogin(ctx context.Context, id uuid.UUID, lastLogin time.Time) error {
	result, err := accountsDB.conn.ExecContext(ctx, "UPDATE accounts SET last_login=$1 WHERE id=$2", lastLogin, id)
	if err != nil {
		return ErrAccountsDB.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return accounts.ErrNoAccount.New("user does not exist")
	}

	return ErrAccountsDB.Wrap(err)
}

// UpdateNonce updates nonce by user.
func (accountsDB *accountsDB) UpdateNonce(ctx context.Context, id uuid.UUID, nonce []byte) error {
	result, err := accountsDB.conn.ExecContext(ctx, "UPDATE accounts SET nonce=$1 WHERE id=$2", nonce, id)
	if err != nil {
		return ErrAccountsDB.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return accounts.ErrNoAccount.New("user does not exist")
	}

	return ErrAccountsDB.Wrap(err)
}

// UpdateStatus updates user status in the database.
func (accountsDB *accountsDB) UpdateStatus(ctx context.Context, id uuid.UUID, status accounts.Status) error {
	result, err := accountsDB.conn.ExecContext(ctx, "UPDATE accounts SET status=$1 WHERE id=$2", status, id)
	if err != nil {
		return ErrAccountsDB.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return accounts.ErrNoAccount.New("user does not exist")
	}

	return ErrAccountsDB.Wrap(err)
}

// Delete deletes a user in the database.
func (accountsDB *accountsDB) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := accountsDB.conn.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", id)
	if err != nil {
		return ErrAccountsDB.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return accounts.ErrNoAccount.New("user does not exist")
	}

	return ErrAccountsDB.Wrap(err)
}
