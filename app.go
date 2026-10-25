package generator

import (
	"context"
	"errors"
	"net"
	"net/mail"

	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"

	"generator/accounts"
	"generator/accounts/auth"
	"generator/collections"
	"generator/collections/images"
	"generator/collections/layers"
	"generator/collections/layers/filestore"
	"generator/collections/metadata"
	"generator/dashboard"
	"generator/internal/filestore/localfilestore"
	"generator/internal/logger"
	"generator/nfts"
	pkgauth "generator/pkg/auth"
	"generator/pkg/chain"
	mail2 "generator/pkg/mail"
	"generator/server/emails"
	"generator/server/httpserver"
	"generator/whitelist"
)

// DB provides access to all databases and database related functionality.
//
// architecture: Master Database.
type DB interface {
	// Close closes underlying db connection.
	Close() error

	// Accounts provides access to accounts db.
	Accounts() accounts.DB

	// Collections provides access to collections db.
	Collections() collections.DB

	// Layers provides access to layers db.
	Layers() layers.DB

	// WhiteList provides access to whitelist db.
	WhiteList() whitelist.DB

	// CreateSchema create tables.
	CreateSchema(ctx context.Context) error
}

// Config is the global configuration for generator.
type Config struct {
	FileStore struct {
		filestore.Config
	} `json:"fileStore"`

	LocalFileStore struct {
		localfilestore.Config
	} `json:"localFileStore"`

	Server struct {
		HTTP   httpserver.Config `json:"server"`
		Emails emails.Config     `json:"emails"`
	} `json:"console"`

	Accounts struct {
		Auth struct {
			TokenSecret string `json:"tokenSecret"`
		} `json:"auth"`
	} `json:"users"`

	NFTs struct {
		nfts.Config
	} `json:"nfts"`

	Dashboard struct {
		dashboard.Config
	} `json:"dashboard"`

	Images struct {
		images.Config
	} `json:"images"`

	Metadata struct {
		metadata.Config
	} `json:"metadata"`

	Collections struct {
		collections.Config
	} `json:"projects"`

	Whitelist struct {
		whitelist.Config
	} `json:"whitelist"`

	Layers struct {
		layers.Config
	} `json:"layers"`
}

// App is the representation of a generator.
type App struct {
	Config   Config
	Log      logger.Logger
	Database DB
	Sender   mail2.Mailer

	// exposes accounts related logic.
	Accounts struct {
		Service *accounts.Service
		Auth    *auth.Service
	}

	// exposes collections related logic.
	Collections struct {
		Service *collections.Service
	}

	// exposes layers related logic.
	Layers struct {
		Service *layers.Service
	}

	// exposes whitelist related logic.
	Whitelist struct {
		Service *whitelist.Service
	}

	// exposes nfts related logic.
	NFTs struct {
		Service *nfts.Service
	}

	// exposes dashboard related logic.
	Dashboard struct {
		Chain   *chain.Reader
		Service *dashboard.Service
	}

	// exposes images relates logic.
	Images struct {
		Service *images.Service
	}

	// exposes metadata relates logic.
	Metadata struct {
		Service *metadata.Service
	}

	// Server web server with web UI.
	Server struct {
		Listener     net.Listener
		Endpoint     *httpserver.Server
		EmailService *emails.Service
	}
}

// New is a constructor for generator.App.
func New(logger logger.Logger, config Config, db DB) (app *App, err error) {
	app = &App{
		Log:      logger,
		Database: db,
	}

	{ // emails setup.
		var sender mail2.Mailer
		if config.Server.Emails.Provider == "mock" {
			sender = &mail2.MockMailer{}
		} else {
			from, err := mail.ParseAddress(config.Server.Emails.From)
			if err != nil {
				logger.Error("invalid email address", errs.Wrap(err))
				return nil, err
			}

			sender = &mail2.SMTPMailer{
				From: *from,
				Auth: mail2.LoginAuth{
					Password: config.Server.Emails.PlainPassword,
					Username: config.Server.Emails.PlainLogin,
				},
				ServerAddress: config.Server.Emails.SMTPServerAddress,
			}
		}

		mailService := emails.NewService(app.Log, sender, config.Server.Emails)
		app.Server.EmailService = mailService
	}

	{ // accounts setup.
		app.Accounts.Service = accounts.NewService(
			app.Database.Accounts(),
		)
		app.Accounts.Auth = auth.NewService(
			app.Database.Accounts(),
			app.Accounts.Service,
			pkgauth.JWTSigner{
				Secret: []byte(config.Accounts.Auth.TokenSecret),
			},
			app.Server.EmailService,
			logger)
	}

	{ // images setup.
		localFileStore := localfilestore.NewFileStore(
			config.LocalFileStore.Config,
		)

		app.Images.Service = images.NewService(
			config.Images.Config,
			localFileStore,
		)
	}

	{ // metadata setup.
		localFileStore := localfilestore.NewFileStore(
			config.LocalFileStore.Config,
		)

		app.Metadata.Service = metadata.NewService(
			config.Metadata.Config,
			localFileStore,
		)
	}

	{ // layers setup.
		fileStore := filestore.NewStore(
			config.FileStore.Config,
		)

		localFileStore := localfilestore.NewFileStore(
			config.LocalFileStore.Config,
		)

		app.Layers.Service = layers.NewService(
			logger,
			config.Layers.Config,
			app.Database.Layers(),
			app.Metadata.Service,
			fileStore,
			localFileStore,
		)
	}

	{ // collections setup.
		localFileStore := localfilestore.NewFileStore(
			config.LocalFileStore.Config,
		)

		app.Collections.Service = collections.NewService(
			config.Collections.Config,
			app.Database.Collections(),
			app.Layers.Service,
			localFileStore,
			app.Images.Service,
			app.Metadata.Service,
		)
	}

	{ // whitelist setup.
		app.Whitelist.Service = whitelist.NewService(
			config.Whitelist.Config,
			app.Database.WhiteList(),
			app.Collections.Service,
		)
	}

	{ // nfts setup.
		app.NFTs.Service = nfts.NewService(
			config.NFTs.Config,
			app.Collections.Service,
			app.Whitelist.Service,
		)
	}

	{ // dashboard setup.
		rpcURL := config.Dashboard.RPCURL
		if rpcURL == "" {
			rpcURL = config.Server.HTTP.Chain.RPCURL
		}
		app.Dashboard.Chain, err = chain.NewReader(rpcURL, config.Dashboard.LogsBlockRange)
		if err != nil {
			return nil, err
		}

		app.Dashboard.Service = dashboard.NewService(config.Dashboard.Config, app.Dashboard.Chain, app.Collections.Service, logger)
	}

	{ // server setup.
		app.Server.Listener, err = net.Listen("tcp", config.Server.HTTP.Address)
		if err != nil {
			return nil, err
		}

		app.Server.Endpoint = httpserver.NewServer(
			config.Server.HTTP,
			logger,
			app.Server.Listener,
			app.Accounts.Auth,
			app.Accounts.Service,
			app.Collections.Service,
			app.Layers.Service,
			app.NFTs.Service,
			app.Whitelist.Service,
			app.Dashboard.Service,
			app.Images.Service,
			app.Metadata.Service,
		)
	}

	return app, nil
}

// Run runs generator.App until it's either closed or it errors.
func (app *App) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// start generator server as a separate goroutine.
	group.Go(func() error {
		return ignoreCancel(app.Server.Endpoint.Run(ctx))
	})

	return group.Wait()
}

// Close closes all the resources.
func (app *App) Close() error {
	var errlist errs.Group

	errlist.Add(app.Server.Endpoint.Close())
	app.Dashboard.Chain.Close()

	return errlist.Err()
}

// we ignore cancellation and stopping errors since they are expected.
func ignoreCancel(err error) error {
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}
