package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/pkg/pagination"
	"generator/whitelist"
)

// ensures that whitelistDB implements whitelist.DB.
var _ whitelist.DB = (*whitelistDB)(nil)

// ErrWhiteList indicates that there was an error in the database.
var ErrWhiteList = errs.Class("whitelist repository error")

// whitelistDB provides access to whitelist db.
//
// architecture: Database
type whitelistDB struct {
	conn *sql.DB
}

// Add adds wallet address to whitelist in the database.
func (whitelistDB *whitelistDB) Add(ctx context.Context, projectID uuid.UUID, item whitelist.Item) error {
	query := `INSERT INTO whitelist(project_id, name, wallet_address) VALUES($1,$2,$3)`
	_, err := whitelistDB.conn.ExecContext(ctx, query, projectID, item.Name, item.WalletAddress)
	return ErrWhiteList.Wrap(err)
}

// AddMany adds the items to whitelist in one transaction: either all of them or none.
func (whitelistDB *whitelistDB) AddMany(ctx context.Context, projectID uuid.UUID, items []whitelist.Item) error {
	query := `INSERT INTO whitelist(project_id, name, wallet_address) VALUES($1,$2,$3)`

	tx, err := whitelistDB.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrWhiteList.Wrap(err)
	}

	for _, item := range items {
		if _, err = tx.ExecContext(ctx, query, projectID, item.Name, item.WalletAddress); err != nil {
			return errs.Combine(ErrWhiteList.Wrap(err), tx.Rollback())
		}
	}

	return ErrWhiteList.Wrap(tx.Commit())
}

// Get returns item of whitelist by project id and wallet address from the database.
func (whitelistDB *whitelistDB) Get(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) (whitelist.Item, error) {
	var item whitelist.Item
	query := "SELECT name FROM whitelist WHERE project_id = $1 AND wallet_address = $2"

	err := whitelistDB.conn.QueryRowContext(ctx, query, projectID, walletAddress).Scan(&item.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return item, whitelist.ErrNoItem.Wrap(err)
	}
	item.WalletAddress = walletAddress

	return item, ErrWhiteList.Wrap(err)
}

// ListByProjectID returns wallet addresses by project id from the database.
func (whitelistDB *whitelistDB) ListByProjectID(ctx context.Context, projectID uuid.UUID) (_ []whitelist.Item, err error) {
	var items []whitelist.Item
	// the order of the leaves defines the Merkle root, so it must not depend on how rows are stored.
	query := `SELECT name, wallet_address FROM whitelist WHERE project_id = $1 ORDER BY wallet_address`

	rows, err := whitelistDB.conn.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, ErrWhiteList.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var item whitelist.Item
		err := rows.Scan(&item.Name, &item.WalletAddress)
		if err != nil {
			return nil, ErrWhiteList.Wrap(err)
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrWhiteList.Wrap(err)
	}

	return items, nil
}

// ListByProjectIDWithPagination returns names and wallet addresses by project id with pagination from the database.
func (whitelistDB *whitelistDB) ListByProjectIDWithPagination(ctx context.Context, projectID uuid.UUID, cursor pagination.Cursor) (_ whitelist.Page, err error) {
	var (
		items         []whitelist.Item
		itemsListPage whitelist.Page
	)

	offset := (cursor.Page - 1) * cursor.Limit
	query := "SELECT name, wallet_address FROM whitelist WHERE project_id = $1 ORDER BY name, wallet_address LIMIT $2 OFFSET $3"

	rows, err := whitelistDB.conn.QueryContext(ctx, query, projectID, cursor.Limit, offset)
	if err != nil {
		return itemsListPage, ErrWhiteList.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var item whitelist.Item
		err := rows.Scan(&item.Name, &item.WalletAddress)
		if err != nil {
			return itemsListPage, ErrWhiteList.Wrap(err)
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return itemsListPage, ErrWhiteList.Wrap(err)
	}

	totalCount, err := whitelistDB.TotalCount(ctx, projectID)
	if err != nil {
		return itemsListPage, ErrWhiteList.Wrap(err)
	}

	pageCount := totalCount / cursor.Limit
	if totalCount%cursor.Limit != 0 {
		pageCount++
	}

	itemsListPage = whitelist.Page{
		Items: items,
		Page: pagination.Page{
			Offset:      offset,
			Limit:       cursor.Limit,
			CurrentPage: cursor.Page,
			PageCount:   pageCount,
			TotalCount:  totalCount,
		},
	}

	return itemsListPage, ErrWhiteList.Wrap(err)
}

// TotalCount counts items of whitelist in the table.
func (whitelistDB *whitelistDB) TotalCount(ctx context.Context, projectID uuid.UUID) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM whitelist WHERE project_id = $1`
	err := whitelistDB.conn.QueryRowContext(ctx, query, projectID).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, whitelist.ErrNoItem.Wrap(err)
	}
	return count, ErrWhiteList.Wrap(err)
}

// ListByProjectIDAndItemName returns items by project id and item name from database.
func (whitelistDB *whitelistDB) ListByProjectIDAndItemName(ctx context.Context, projectID uuid.UUID, name string, cursor pagination.Cursor) (_ whitelist.Page, err error) {
	var (
		whitelistPage whitelist.Page
		values        []interface{}
		totalCount    int
	)

	values = append(values, name, name+" %", "% "+name, "% "+name+" %", projectID.String())

	extendedValues := make([]interface{}, len(values))
	copy(extendedValues, values)
	extendedValues = append(extendedValues, cursor.Limit)
	offset := (cursor.Page - 1) * cursor.Limit
	extendedValues = append(extendedValues, offset)

	query :=
		`SELECT 
			name, wallet_address 
		FROM 
			whitelist 
		WHERE
			(name LIKE $1 OR name LIKE $2 OR name LIKE $3 OR name LIKE $4) AND project_id = $5
		ORDER BY name, wallet_address
		LIMIT $6
		OFFSET $7`

	rows, err := whitelistDB.conn.QueryContext(ctx, query, extendedValues...)
	if err != nil {
		return whitelistPage, ErrWhiteList.Wrap(err)
	}
	defer func() {
		err = ErrWhiteList.Wrap(errs.Combine(err, rows.Close()))
	}()

	var items []whitelist.Item
	for rows.Next() {
		item := whitelist.Item{}
		if err = rows.Scan(&item.Name, &item.WalletAddress); err != nil {
			return whitelistPage, ErrWhiteList.Wrap(err)
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return whitelistPage, ErrWhiteList.Wrap(err)
	}

	query = "SELECT COUNT(*) FROM whitelist WHERE (name LIKE $1 OR name LIKE $2 OR name LIKE $3 OR name LIKE $4) AND project_id = $5"
	err = whitelistDB.conn.QueryRowContext(ctx, query, values...).Scan(&totalCount)
	if errors.Is(err, sql.ErrNoRows) {
		return whitelistPage, whitelist.ErrNoItem.Wrap(err)
	}

	pageCount := totalCount / cursor.Limit
	if totalCount%cursor.Limit != 0 {
		pageCount++
	}

	whitelistPage = whitelist.Page{
		Items: items,
		Page: pagination.Page{
			Offset:      offset,
			Limit:       cursor.Limit,
			CurrentPage: cursor.Page,
			PageCount:   pageCount,
			TotalCount:  totalCount,
		},
	}

	return whitelistPage, ErrWhiteList.Wrap(err)
}

// Delete deletes wallet address from whitelist in the database.
func (whitelistDB *whitelistDB) Delete(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) error {
	query := `DELETE FROM whitelist WHERE project_id = $1 AND wallet_address = $2`

	result, err := whitelistDB.conn.ExecContext(ctx, query, projectID, walletAddress)
	if err != nil {
		return ErrWhiteList.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return whitelist.ErrNoItem.New("")
	}

	return ErrWhiteList.Wrap(err)
}
