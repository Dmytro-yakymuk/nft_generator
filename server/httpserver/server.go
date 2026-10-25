package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"

	"generator/accounts"
	"generator/accounts/auth"
	"generator/collections"
	"generator/collections/images"
	"generator/collections/layers"
	"generator/collections/metadata"
	"generator/dashboard"
	"generator/internal/logger"
	"generator/nfts"
	pkgauth "generator/pkg/auth"
	"generator/pkg/chain"
	"generator/server/httpserver/handlers"
	"generator/whitelist"
)

var (
	// Error is an error class that indicates internal http server error.
	Error = errs.Class("console web server error")
)

// Config contains configuration for console web server.
type Config struct {
	Address         string `json:"address"`
	StaticDir       string `json:"staticDir"`
	LayersImagesDir string `json:"layersImagesDir"`
	ProjectsDir     string `json:"projectsDir"`

	Auth struct {
		CookieName string `json:"cookieName"`
		Path       string `json:"path"`
	} `json:"auth"`

	// Chain is passed to the web app and mint sites, which switch MetaMask to this network.
	Chain chain.Config `json:"chain"`
	// IPFSGatewayURL is the HTTP gateway the web app links IPFS files through, e.g. http://localhost:8080.
	IPFSGatewayURL string `json:"ipfsGatewayUrl"`
	// MintSiteDomain, if set, serves each published mint site at its own subdomain <slug>.<MintSiteDomain>
	// besides /sites/<slug>; it needs a wildcard DNS record (and certificate) for *.<MintSiteDomain>,
	// while *.localhost resolves to this machine without one.
	MintSiteDomain string `json:"mintSiteDomain"`
}

// appConfig is the runtime configuration of the web app, injected into index.html.
type appConfig struct {
	Chain          chain.Config `json:"chain"`
	IPFSGatewayURL string       `json:"ipfsGatewayUrl"`
	// IPFSPinning names the service published collections can be pinned in, empty if there is none.
	IPFSPinning string `json:"ipfsPinning"`
	// MintSiteDomain is the domain whose subdomains serve the published mint sites, empty if there is none.
	MintSiteDomain string `json:"mintSiteDomain"`
}

// Server represents console web server.
//
// architecture: Endpoint
type Server struct {
	log    logger.Logger
	config Config
	// ipfsPinning is the pinning service passed to the web app.
	ipfsPinning string

	listener net.Listener
	server   http.Server

	authService *auth.Service
	cookieAuth  *pkgauth.SessionCookie

	templates struct {
		index    *template.Template
		auth     *handlers.AuthTemplates
		projects *template.Template
	}
}

// NewServer is a constructor for console web server.
func NewServer(config Config, log logger.Logger, listener net.Listener, userAuth *auth.Service, users *accounts.Service,
	projects *collections.Service, layers *layers.Service, nfts *nfts.Service, whitelist *whitelist.Service, dashboard *dashboard.Service,
	images *images.Service, metadata *metadata.Service) *Server {
	server := &Server{
		log:         log,
		config:      config,
		ipfsPinning: projects.PinningService(),
		listener:    listener,
		authService: userAuth,
		cookieAuth: pkgauth.NewSessionCookie(pkgauth.SessionSettings{
			Name: config.Auth.CookieName,
			Path: config.Auth.Path,
		}),
	}

	err := server.initializeTemplates()
	if err != nil {
		return nil
	}

	projectsController := handlers.NewProjects(log, projects, layers, images, metadata, server.templates.projects, config.Chain)
	usersController := handlers.NewUsers(log, users)
	layersController := handlers.NewLayers(log, layers)
	contractController := handlers.NewContracts(log, nfts)
	whitelistController := handlers.NewWhitelist(log, whitelist, projects)
	authController := handlers.NewAuth(server.log, server.authService, server.cookieAuth, server.templates.auth)
	dashboardController := handlers.NewDashboard(log, dashboard)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v0").Subrouter()

	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	metamaskRouter := authRouter.PathPrefix("/metamask").Subrouter()
	metamaskRouter.HandleFunc("/register", authController.MetamaskRegister).Methods(http.MethodPost)
	metamaskRouter.HandleFunc("/nonce", authController.Nonce).Methods(http.MethodGet)
	metamaskRouter.HandleFunc("/login", authController.MetamaskLogin).Methods(http.MethodPost)
	authRouter.HandleFunc("/logout", authController.Logout).Methods(http.MethodPost)

	accountsRouter := apiRouter.PathPrefix("/accounts").Subrouter()
	accountsRouter.Use(server.withAuth)
	accountsRouter.HandleFunc("", usersController.Get).Methods(http.MethodGet)
	accountsRouter.HandleFunc("", usersController.Update).Methods(http.MethodPatch)

	// Routes with a collection in the path ({id} or {projectId}) are available only to its owner.
	ownerOnly := withCollectionOwner(log, projects)

	collectionsRouter := apiRouter.PathPrefix("/collections").Subrouter()
	collectionsRouter.Use(server.withAuth)
	collectionsRouter.Use(ownerOnly)
	collectionsRouter.HandleFunc("", projectsController.CreateProject).Methods(http.MethodPost)
	collectionsRouter.HandleFunc("", projectsController.List).Methods(http.MethodGet)
	collectionsRouter.HandleFunc("/{id}", projectsController.Get).Methods(http.MethodGet)
	collectionsRouter.HandleFunc("/{id}/generate-images", projectsController.Generate).Methods(http.MethodPost)
	collectionsRouter.HandleFunc("/{id}/abort-generation", projectsController.AbortGeneration).Methods(http.MethodPut)
	collectionsRouter.HandleFunc("/{id}/generate-metadata", projectsController.GenerateMetadata).Methods(http.MethodPost)
	collectionsRouter.HandleFunc("/{id}/generated-images-count", projectsController.GeneratedImagesCount).Methods(http.MethodGet)
	collectionsRouter.HandleFunc("/{id}/upload-files", projectsController.UploadFiles).Methods(http.MethodPost)
	collectionsRouter.HandleFunc("/{id}/ipfs", projectsController.PublishToIPFS).Methods(http.MethodPost)
	collectionsRouter.HandleFunc("/{id}/ipfs/pin", projectsController.PinIPFS).Methods(http.MethodPost)
	collectionsRouter.HandleFunc("/{id}/nfts", projectsController.DeleteNFTs).Methods(http.MethodDelete)
	collectionsRouter.HandleFunc("/{id}", projectsController.Delete).Methods(http.MethodDelete)
	collectionsRouter.HandleFunc("/{id}", projectsController.Update).Methods(http.MethodPatch)
	collectionsRouter.HandleFunc("/{id}/contractAddress", projectsController.UpdateContractAddress).Methods(http.MethodPatch)
	collectionsRouter.HandleFunc("/{id}/website", projectsController.PublishSite).Methods(http.MethodPut)
	collectionsRouter.HandleFunc("/{id}/website", projectsController.UnpublishSite).Methods(http.MethodDelete)

	mintSiteRouter := collectionsRouter.PathPrefix("/{projectId}/mint-site").Subrouter()
	mintSiteRouter.Use(server.withAuth)
	mintSiteRouter.HandleFunc("", projectsController.GetProjectPage).Methods(http.MethodGet)
	mintSiteRouter.HandleFunc("", projectsController.DeleteProjectSite).Methods(http.MethodDelete)

	logoRouter := mintSiteRouter.PathPrefix("/logo").Subrouter()
	logoRouter.Use(server.withAuth)
	logoRouter.HandleFunc("", projectsController.UpdateLogo).Methods(http.MethodPut)
	logoRouter.HandleFunc("", projectsController.DeleteLogo).Methods(http.MethodDelete)

	headerSectionRouter := mintSiteRouter.PathPrefix("/header-section").Subrouter()
	headerSectionRouter.Use(server.withAuth)
	headerSectionRouter.HandleFunc("", projectsController.UpdateHeaderSection).Methods(http.MethodPatch)
	headerSectionRouter.HandleFunc("", projectsController.UpdateHeaderSectionImage).Methods(http.MethodPut)
	headerSectionRouter.HandleFunc("", projectsController.GetHeaderSection).Methods(http.MethodGet)
	headerSectionRouter.HandleFunc("/image", projectsController.DeleteHeaderSectionImage).Methods(http.MethodDelete)

	aboutSectionRouter := mintSiteRouter.PathPrefix("/about-section").Subrouter()
	aboutSectionRouter.Use(server.withAuth)
	aboutSectionRouter.HandleFunc("", projectsController.UpdateAboutSection).Methods(http.MethodPatch)
	aboutSectionRouter.HandleFunc("", projectsController.UpdateAboutSectionImage).Methods(http.MethodPut)
	aboutSectionRouter.HandleFunc("", projectsController.GetAboutSection).Methods(http.MethodGet)
	aboutSectionRouter.HandleFunc("/image", projectsController.DeleteAboutSectionImage).Methods(http.MethodDelete)

	membersRouter := mintSiteRouter.PathPrefix("/members").Subrouter()
	membersRouter.Use(server.withAuth)
	membersRouter.HandleFunc("", projectsController.AddMember).Methods(http.MethodPost)
	membersRouter.HandleFunc("/{memberId}", projectsController.UpdateMemberImage).Methods(http.MethodPut)
	membersRouter.HandleFunc("", projectsController.ListMembers).Methods(http.MethodGet)
	membersRouter.HandleFunc("/{memberId}", projectsController.UpdateMember).Methods(http.MethodPatch)
	membersRouter.HandleFunc("/{memberId}", projectsController.DeleteMember).Methods(http.MethodDelete)
	membersRouter.HandleFunc("/{memberId}/image", projectsController.DeleteMemberImage).Methods(http.MethodDelete)

	milestonesRouter := mintSiteRouter.PathPrefix("/milestones").Subrouter()
	milestonesRouter.Use(server.withAuth)
	milestonesRouter.HandleFunc("", projectsController.AddMilestone).Methods(http.MethodPost)
	milestonesRouter.HandleFunc("", projectsController.ListMilestones).Methods(http.MethodGet)
	milestonesRouter.HandleFunc("/{milestoneId}", projectsController.UpdateMilestone).Methods(http.MethodPatch)
	milestonesRouter.HandleFunc("/{milestoneId}", projectsController.DeleteMilestone).Methods(http.MethodDelete)

	contentRouter := mintSiteRouter.PathPrefix("/content").Subrouter()
	contentRouter.Use(server.withAuth)
	contentRouter.HandleFunc("", projectsController.UpdateContentImage).Methods(http.MethodPut)
	contentRouter.HandleFunc("", projectsController.ListContentImages).Methods(http.MethodGet)
	contentRouter.HandleFunc("/{name}", projectsController.DeleteContentImage).Methods(http.MethodDelete)

	socialNetworksRouter := mintSiteRouter.PathPrefix("/social-networks").Subrouter()
	socialNetworksRouter.Use(server.withAuth)
	socialNetworksRouter.HandleFunc("", projectsController.AddSocialNetwork).Methods(http.MethodPost)
	socialNetworksRouter.HandleFunc("", projectsController.ListSocialNetworks).Methods(http.MethodGet)
	socialNetworksRouter.HandleFunc("/{socialNetworkId}", projectsController.UpdateSocialNetwork).Methods(http.MethodPatch)
	socialNetworksRouter.HandleFunc("/{socialNetworkId}", projectsController.DeleteSocialNetwork).Methods(http.MethodDelete)

	layersRouter := collectionsRouter.PathPrefix("/{projectId}/layers").Subrouter()
	layersRouter.Use(server.withAuth)
	layersRouter.HandleFunc("", layersController.Create).Methods(http.MethodPost)
	layersRouter.HandleFunc("", layersController.List).Methods(http.MethodGet)
	layersRouter.HandleFunc("/{layerId}", layersController.ListLayerImageNames).Methods(http.MethodGet)
	layersRouter.HandleFunc("/{layerId}", layersController.UpdatePriority).Methods(http.MethodPatch)
	layersRouter.HandleFunc("/{layerId}", layersController.Delete).Methods(http.MethodDelete)
	layersRouter.HandleFunc("/{layerId}/weights", layersController.UpdateWeights).Methods(http.MethodPut)

	layerImagesRouter := layersRouter.PathPrefix("/{layerId}/images").Subrouter()
	layerImagesRouter.Use(server.withAuth)
	layerImagesRouter.HandleFunc("", layersController.CreateLayerImage).Methods(http.MethodPost)
	layerImagesRouter.HandleFunc("/{name}", layersController.UpdateImageName).Methods(http.MethodPatch)
	layerImagesRouter.HandleFunc("/{name}", layersController.DeleteLayerImage).Methods(http.MethodDelete)

	contractRouter := apiRouter.PathPrefix("/contracts").Subrouter()
	contractRouter.Use(server.withAuth)
	contractRouter.Use(ownerOnly)
	contractRouter.HandleFunc("/data/{projectId}", contractController.GetDeployContractTransactionData).Methods(http.MethodGet)

	whitelistRouter := apiRouter.PathPrefix("/whitelist").Subrouter()
	whitelistRouter.Use(server.withAuth)
	whitelistRouter.Use(ownerOnly)
	whitelistRouter.HandleFunc("", whitelistController.Add).Methods(http.MethodPost)
	whitelistRouter.HandleFunc("/{projectId}/merkle-root", whitelistController.GetMerkleRoot).Methods(http.MethodGet)
	whitelistRouter.HandleFunc("/{projectId}/excel-file", whitelistController.AddFromExcelFile).Methods(http.MethodPost)
	whitelistRouter.HandleFunc("/{projectId}", whitelistController.List).Methods(http.MethodGet)
	whitelistRouter.HandleFunc("/{projectId}/{walletAddress}", whitelistController.Delete).Methods(http.MethodDelete)
	whitelistRouter.HandleFunc("/{projectId}/{walletAddress}/merkle-proof", whitelistController.GetMerkleProof).Methods(http.MethodGet)

	dashboardRouter := apiRouter.PathPrefix("/dashboard/{projectId}").Subrouter()
	dashboardRouter.Use(server.withAuth)
	dashboardRouter.HandleFunc("/events", dashboardController.Events).Methods(http.MethodGet)
	dashboardRouter.HandleFunc("/stats", dashboardController.CollectionStats).Methods(http.MethodGet)
	dashboardRouter.HandleFunc("/sales", dashboardController.EventData).Methods(http.MethodGet)

	ls := http.FileServer(http.Dir(server.config.LayersImagesDir))
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", ls))

	router.HandleFunc("/metadata/{id}/{tokenId:[0-9]+}", projectsController.TokenMetadata).Methods(http.MethodGet)

	// Published mint sites and the data their scripts need are available without authentication.
	router.HandleFunc("/sites/{slug}", projectsController.PublicSitePage).Methods(http.MethodGet)
	handleMintSiteHosts(router, server.config.MintSiteDomain, projectsController.PublicSitePage)
	publicRouter := apiRouter.PathPrefix("/public").Subrouter()
	publicRouter.HandleFunc("/whitelist/{projectId}/{walletAddress}/merkle-proof", whitelistController.GetPublicMerkleProof).Methods(http.MethodGet)

	css := http.FileServer(http.Dir(server.config.ProjectsDir))
	router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css", css))

	js := http.FileServer(http.Dir(server.config.ProjectsDir))
	router.PathPrefix("/static/js/").Handler(http.StripPrefix("/static/js", js))

	projectImages := http.FileServer(http.Dir(server.config.ProjectsDir))
	router.PathPrefix("/static/images/").Handler(http.StripPrefix("/static", projectImages))

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", newPrecompressed(server.config.StaticDir)))
	router.PathPrefix("/").HandlerFunc(server.appHandler)

	server.server = http.Server{
		Handler: router,
	}

	return server
}

// Run starts the server that host webapp and api endpoint.
func (server *Server) Run(ctx context.Context) (err error) {
	err = server.initializeTemplates()
	if err != nil {
		return Error.Wrap(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return Error.Wrap(server.server.Shutdown(context.Background()))
	})
	group.Go(func() error {
		defer cancel()
		err := server.server.Serve(server.listener)
		isCancelled := errs.IsFunc(err, func(err error) bool { return errors.Is(err, context.Canceled) })
		if isCancelled || errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		return Error.Wrap(err)
	})

	return Error.Wrap(group.Wait())
}

// Close closes server and underlying listener.
func (server *Server) Close() error {
	return Error.Wrap(server.server.Close())
}

// appHandler is web app http handler function.
func (server *Server) appHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	header.Set("Content-Type", "text/html; charset=UTF-8")
	// header.Set("X-Content-Type-Options", "nosniff").
	header.Set("Referrer-Policy", "same-origin")

	if server.templates.index == nil {
		server.log.Error("index template is not set", nil)
		return
	}

	if err := server.templates.index.Execute(w, appConfig{Chain: server.config.Chain, IPFSGatewayURL: server.config.IPFSGatewayURL, IPFSPinning: server.ipfsPinning,
		MintSiteDomain: server.config.MintSiteDomain}); err != nil {
		server.log.Error("index template could not be executed", err)
		return
	}
}

// serveUnauthorized replies 401 with a JSON error, as the API handlers do.
func (server *Server) serveUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "authorization required"}); err != nil {
		server.log.Error("failed to write json error response", Error.Wrap(err))
	}
}

// withAuth performs initial authorization before every request.
func (server *Server) withAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := server.cookieAuth.GetToken(r)
		if err != nil {
			server.serveUnauthorized(w)
			return
		}

		claims, err := server.authService.Authorize(ctx, token)
		if err != nil {
			server.serveUnauthorized(w)
			return
		}

		ctx = pkgauth.SetTokenClaims(ctx, claims)

		handler.ServeHTTP(w, r.Clone(ctx))
	})
}

// initializeTemplates is used to initialize all templates.
func (server *Server) initializeTemplates() (err error) {
	server.templates.index, err = template.ParseFiles(filepath.Join(server.config.StaticDir, "dist", "index.html"))
	if err != nil {
		server.log.Error("dist folder is not generated. use 'npm run build' command", err)
		return err
	}

	server.templates.projects, err = template.ParseFiles(filepath.Join(server.config.ProjectsDir, "projects.html"))
	if err != nil {
		server.log.Error("could not parse html representation of projects mint site", err)
		return err
	}

	return nil
}

// handleMintSiteHosts serves the root page of <slug>.<domain> (any port) as the mint site with that address;
// the scripts and styles of the page come from the same host by the usual /static and /api paths.
// It must be registered before the web app route, which takes every other path.
func handleMintSiteHosts(router *mux.Router, domain string, sitePage http.HandlerFunc) {
	if domain == "" {
		return
	}
	router.Host("{slug:[a-z0-9-]+}." + domain).Path("/").Methods(http.MethodGet).HandlerFunc(sitePage)
}
