package storage

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq" // using postgres driver.
	"github.com/zeebo/errs"

	"generator"
	"generator/accounts"
	"generator/collections"
	"generator/collections/layers"
	"generator/whitelist"
)

// ensures that storage implements generator.DB.
var _ generator.DB = (*storage)(nil)

var (
	// Error is the default generator error class.
	Error = errs.Class("generator storage error")
)

// storage combines access to different database tables with a record
// of the db driver, db implementation, and db source URL.
//
// architecture: Master Database
type storage struct {
	conn *sql.DB
}

// New returns generator.DB postgresql implementation.
func New(databaseURL string) (generator.DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &storage{conn: conn}, nil
}

// CreateSchema create schema for all tables and databases.
func (db *storage) CreateSchema(ctx context.Context) (err error) {
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS accounts (
            id               BYTEA  PRIMARY KEY       NOT NULL,
            name             VARCHAR                  NOT NULL,
            wallet_address   VARCHAR                  NOT NULL,
            nonce            BYTEA                    NOT NUll,
            email            VARCHAR                  NOT NULL,
            email_normalized VARCHAR                  NOT NULL,
            password_hash    BYTEA                    NOT NULL,
            last_login       TIMESTAMP WITH TIME ZONE NOT NULL,
            status           INTEGER                  NOT NULL,
            is_subscribed    BOOLEAN                  NOT NULL,
            created_at       TIMESTAMP WITH TIME ZONE NOT NULL
        );
        CREATE TABLE IF NOT EXISTS collections (
            id                   BYTEA   PRIMARY KEY                               NOT NULL,
            name                 VARCHAR UNIQUE                                    NOT NULL,
            nft_count            INTEGER                                           NOT NULL,
            user_id              BYTEA   REFERENCES accounts(id) ON DELETE CASCADE NOT NULL,
            contract_address     BYTEA                                             NOT NULL,
            are_images_generated BOOLEAN                                           NOT NULL,
            is_website_published BOOLEAN                                           NOT NULL
        );
        CREATE TABLE IF NOT EXISTS layers (
            id         BYTEA   PRIMARY KEY                                  NOT NULL,
            project_id BYTEA   REFERENCES collections(id) ON DELETE CASCADE NOT NULL,
            name       VARCHAR                                              NOT NULL,
            priority   INTEGER                                              NOT NULL
        );
        CREATE TABLE IF NOT EXISTS whitelist (
            project_id     BYTEA REFERENCES collections(id) ON DELETE CASCADE NOT NULL,
            name           VARCHAR                                            NOT NULL,
            wallet_address BYTEA                                              NOT NULL,
            PRIMARY KEY(project_id, wallet_address)
        );
        CREATE TABLE IF NOT EXISTS collection_mint_site (
            project_id         BYTEA   PRIMARY KEY REFERENCES collections(id) ON DELETE CASCADE NOT NULL,
            title_header       VARCHAR,
            description_header VARCHAR,
            title_about        VARCHAR,
            description_about  VARCHAR
        );
        CREATE TABLE IF NOT EXISTS collection_mint_site_members (
            id         BYTEA   PRIMARY KEY                                  NOT NULL,
            project_id BYTEA   REFERENCES collections(id) ON DELETE CASCADE NOT NULL,
            name       VARCHAR                                              NOT NULL,
            role       VARCHAR,
            created_at TIMESTAMP WITH TIME ZONE                             NOT NULL
        );
        CREATE TABLE IF NOT EXISTS collection_mint_site_milestones (
            id          BYTEA     PRIMARY KEY                                  NOT NULL,
            project_id  BYTEA     REFERENCES collections(id) ON DELETE CASCADE NOT NULL,
            title       VARCHAR                                                NOT NULL,
            date        VARCHAR                                                NOT NULL,
            description VARCHAR
        );
        CREATE TABLE IF NOT EXISTS collection_mint_site_social_networks (
            id          BYTEA     PRIMARY KEY                                  NOT NULL,
            project_id  BYTEA     REFERENCES collections(id) ON DELETE CASCADE NOT NULL,
            title       VARCHAR                                                NOT NULL,
            link        VARCHAR                                                NOT NULL
        );
        ALTER TABLE collections ADD COLUMN IF NOT EXISTS website_slug VARCHAR UNIQUE;
        ALTER TABLE collections ADD COLUMN IF NOT EXISTS contract_block BIGINT NOT NULL DEFAULT 0;
        ALTER TABLE collections ADD COLUMN IF NOT EXISTS images_cid VARCHAR NOT NULL DEFAULT '';
        ALTER TABLE collections ADD COLUMN IF NOT EXISTS metadata_cid VARCHAR NOT NULL DEFAULT '';
        ALTER TABLE collections ADD COLUMN IF NOT EXISTS ipfs_pinned VARCHAR NOT NULL DEFAULT '';
        ALTER TABLE layers ADD COLUMN IF NOT EXISTS none_weight INTEGER NOT NULL DEFAULT 0;
        CREATE TABLE IF NOT EXISTS layer_weights (
            layer_id BYTEA   REFERENCES layers(id) ON DELETE CASCADE NOT NULL,
            name     VARCHAR                                        NOT NULL,
            weight   INTEGER                                        NOT NULL CHECK (weight > 0),
            PRIMARY KEY(layer_id, name)
        );`

	_, err = db.conn.ExecContext(ctx, createTableQuery)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}

// Close closes underlying db connection.
func (db *storage) Close() error {
	return Error.Wrap(db.conn.Close())
}

// Accounts provides access to accounts db.
func (db *storage) Accounts() accounts.DB {
	return &accountsDB{conn: db.conn}
}

// Collections provides access to collections db.
func (db *storage) Collections() collections.DB {
	return &collectionsDB{conn: db.conn}
}

// Layers provides access to layers db.
func (db *storage) Layers() layers.DB {
	return &layersDB{conn: db.conn}
}

// WhiteList provides access to whitelist db.
func (db *storage) WhiteList() whitelist.DB {
	return &whitelistDB{conn: db.conn}
}
