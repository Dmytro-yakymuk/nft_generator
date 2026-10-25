package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"generator/collections"
	"generator/collections/images"
	"generator/collections/layers"
	"generator/collections/metadata"
	"generator/internal/filestore/localfilestore"
	"generator/internal/logger"
	"generator/nfts"
	pkgauth "generator/pkg/auth"
	"generator/pkg/chain"
	"generator/pkg/filebase"
	"generator/pkg/ipfs"
)

// ErrProjects is an internal error type for collections controller.
var ErrProjects = errs.Class("collections controller error")

// Projects is a mvc controller that handles all collections related views.
type Projects struct {
	log logger.Logger

	projects *collections.Service
	layers   *layers.Service
	images   *images.Service
	metadata *metadata.Service

	template *template.Template
	chain    chain.Config
}

// NewProjects is a constructor for collections controller.
func NewProjects(log logger.Logger, projects *collections.Service, layers *layers.Service, images *images.Service, metadata *metadata.Service,
	template *template.Template, chain chain.Config) *Projects {
	projectsController := &Projects{
		log:      log,
		projects: projects,
		layers:   layers,
		images:   images,
		metadata: metadata,
		template: template,
		chain:    chain,
	}

	return projectsController
}

// mintSitePage is the data of the mint site template.
type mintSitePage struct {
	collections.Site
	// Config is passed to the page scripts as window.MINT_SITE.
	Config mintSiteConfig
}

// mintSiteConfig describes the collection contract and network for the mint site scripts.
type mintSiteConfig struct {
	ProjectID       uuid.UUID `json:"projectId"`
	ContractAddress string    `json:"contractAddress"`
	NFTCount        int       `json:"nftCount"`
	// the sale terms of contracts deployed before they became constructor parameters;
	// the mint site reads the terms of newer contracts from the contract itself.
	MintPrice  string       `json:"mintPrice"`
	MaxPerMint int          `json:"maxPerMint"`
	Chain      chain.Config `json:"chain"`
}

// CreateProject is an endpoint that creates collections.
func (controller *Projects) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	claimsMap, err := pkgauth.GetTokenClaims(ctx)
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrProjects.Wrap(err))
		return
	}

	userID, err := uuid.Parse(claimsMap["userId"])
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrProjects.Wrap(err))
		return
	}

	project := collections.Collection{
		ID:     uuid.New(),
		UserID: userID,
	}

	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.Create(ctx, project)
	if err != nil {
		controller.log.Error("could not create projects", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// Generate is an endpoint that generates nft for collections.
func (controller *Projects) Generate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.Generate(ctx, projectID)
	if err != nil {
		switch {
		case layers.ErrNoLayerFile.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		case layers.ErrCountLayers.Has(err), layers.ErrGenerationProcess.Has(err), layers.ErrRareCombinations.Has(err),
			errs.Is(err, layers.ErrInvalidDimension):
			controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
			return
		default:
			controller.log.Error("could not generate project nfts", ErrProjects.Wrap(err))
			controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
			return
		}
	}

}

// GenerateMetadata is an endpoint that generates metadata of nfts for collections.
func (controller *Projects) GenerateMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	if err = controller.projects.GenerateMetadata(ctx, projectID); err != nil {
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		if localfilestore.ErrNoFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not generate metadata of nfts for projects", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// PublishToIPFS is an endpoint that publishes the generated collection to IPFS and returns the
// collection with the CIDs of its images and metadata.
func (controller *Projects) PublishToIPFS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	projectID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	project, err := controller.projects.PublishToIPFS(ctx, projectID)
	if err != nil {
		switch {
		case collections.ErrNoCollection.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
		case collections.ErrNoGeneratedFiles.Has(err):
			controller.serveError(w, http.StatusConflict, ErrProjects.Wrap(err))
		case ipfs.Error.Has(err):
			controller.log.Error("could not publish the collection to IPFS", ErrProjects.Wrap(err))
			controller.serveError(w, http.StatusBadGateway, ErrProjects.Wrap(err))
		default:
			controller.log.Error("could not publish the collection to IPFS", ErrProjects.Wrap(err))
			controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		}
		return
	}

	if err = json.NewEncoder(w).Encode(project); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
	}
}

// PinIPFS is an endpoint that pins the published collection in the pinning service and returns the
// collection with the name of the service.
func (controller *Projects) PinIPFS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	projectID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	project, err := controller.projects.PinIPFS(ctx, projectID)
	if err != nil {
		switch {
		case collections.ErrNoCollection.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
		case collections.ErrNotPublished.Has(err), collections.ErrPinningDisabled.Has(err):
			controller.serveError(w, http.StatusConflict, ErrProjects.Wrap(err))
		case ipfs.Error.Has(err), filebase.Error.Has(err):
			controller.log.Error("could not pin the collection", ErrProjects.Wrap(err))
			controller.serveError(w, http.StatusBadGateway, ErrProjects.Wrap(err))
		default:
			controller.log.Error("could not pin the collection", ErrProjects.Wrap(err))
			controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		}
		return
	}

	if err = json.NewEncoder(w).Encode(project); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
	}
}

// UploadFiles is an endpoint that uploads files for collections.
func (controller *Projects) UploadFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	if err = controller.projects.UploadFiles(ctx, projectID); err != nil {
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		if localfilestore.ErrNoFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not upload files for projects", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// List is an endpoint that returns all users collections.
func (controller *Projects) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	claimsMap, err := pkgauth.GetTokenClaims(ctx)
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrProjects.Wrap(err))
		return
	}

	id, err := uuid.Parse(claimsMap["userId"])
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrProjects.Wrap(err))
		return
	}

	userProjects, err := controller.projects.List(ctx, id)
	if err != nil {
		controller.log.Error("could not list user projects", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(userProjects); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// Get is an endpoint that returns collection by id.
func (controller *Projects) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	project, err := controller.projects.Get(ctx, id)
	if err != nil {
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not get project", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(project); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// GeneratedImagesCount is an endpoint that count of generated images.
func (controller *Projects) GeneratedImagesCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	count, err := controller.images.GeneratedImagesCount(ctx, id)
	if err != nil {
		if localfilestore.ErrNoFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not count generated images", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(count); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// TokenMetadata is a public endpoint that returns ERC-721 metadata of the token,
// so the contract base URI can point to the server.
func (controller *Projects) TokenMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	tokenID, err := strconv.Atoi(vars["tokenId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	metadata, err := controller.projects.TokenMetadata(ctx, id, tokenID)
	if err != nil {
		if collections.ErrNoCollection.Has(err) || collections.ErrNoToken.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not get token metadata", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(metadata); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// DeleteNFTs is an endpoint that deletes generated images and metadata files.
func (controller *Projects) DeleteNFTs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	if err = controller.images.DeleteGeneratedImages(ctx, id); err != nil {
		if collections.ErrNoGeneratedFiles.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}

		controller.log.Error("could not delete generated images", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = controller.metadata.Delete(ctx, id); err != nil {
		if collections.ErrNoGeneratedFiles.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}

		controller.log.Error("could not delete metadata files", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// AbortGeneration is an endpoint that aborts generation.
func (controller *Projects) AbortGeneration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.layers.Abort(id)
	if errs.Is(err, layers.ErrNoGenerationProcess) {
		controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
		return
	}
}

// Update is an endpoints that updates user collections.
func (controller *Projects) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	var update collections.Update
	if err = json.NewDecoder(r.Body).Decode(&update); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.Update(ctx, id, update)
	if err != nil {
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not update project", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateContractAddress is an endpoints that saves the deployed contract address and its deployment block for the collection.
func (controller *Projects) UpdateContractAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id, userID, ok := controller.projectAndUser(w, r)
	if !ok {
		return
	}

	request := struct {
		ContractAddress common.Address `json:"contractAddress"`
		ContractBlock   uint64         `json:"contractBlock"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	if err := controller.projects.UpdateContractAddress(ctx, id, userID, request.ContractAddress, request.ContractBlock); err != nil {
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		if collections.ErrCollectionOwnership.Has(err) {
			controller.serveError(w, http.StatusForbidden, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not update contract address for the project", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// Delete is an endpoints that deletes user collections.
func (controller *Projects) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.Delete(ctx, id)
	if err != nil {
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not delete project", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateHeaderSection is an endpoint that updates header section for collection mint site.
func (controller *Projects) UpdateHeaderSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	var headerSection collections.Section
	if err = json.NewDecoder(r.Body).Decode(&headerSection); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}
	headerSection.ProjectID = id

	err = controller.projects.UpdateHeaderSection(ctx, headerSection)
	if err != nil {
		controller.log.Error("could not update header section", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateHeaderSectionImage is an endpoint that updates image for header section.
func (controller *Projects) UpdateHeaderSectionImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Image string `json:"image"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	img, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.New("invalid image"))
		return
	}

	err = controller.projects.UpdateHeaderSectionImage(ctx, id, bytes.NewBuffer(img))
	if err != nil {
		controller.log.Error("could not update header section image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// GetHeaderSection is an endpoint that returns header section of collection mint site.
func (controller *Projects) GetHeaderSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	headerSection, err := controller.projects.GetHeaderSection(ctx, id)
	if err != nil {
		if collections.ErrNoHeaderSection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not get header section", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(headerSection); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// DeleteHeaderSectionImage is an endpoint that deletes header section image.
func (controller *Projects) DeleteHeaderSectionImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteHeaderSectionImage(ctx, id)
	if err != nil {
		controller.log.Error("could not delete header section image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateLogo is an endpoint that updates collection logo.
func (controller *Projects) UpdateLogo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Image string `json:"image"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	img, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.New("invalid image"))
		return
	}

	err = controller.projects.UpdateLogo(ctx, id, bytes.NewBuffer(img))
	if err != nil {
		controller.log.Error("could not update project logo", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// DeleteLogo is an endpoint that deletes collection logo.
func (controller *Projects) DeleteLogo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteLogo(ctx, id)
	if err != nil {
		controller.log.Error("could not delete project logo", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateAboutSection is an endpoint that updates about section for collection mint site.
func (controller *Projects) UpdateAboutSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	var aboutSection collections.Section
	if err = json.NewDecoder(r.Body).Decode(&aboutSection); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}
	aboutSection.ProjectID = id

	err = controller.projects.UpdateAboutSection(ctx, aboutSection)
	if err != nil {
		controller.log.Error("could not create about section", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateAboutSectionImage is an endpoint that updates image for about section.
func (controller *Projects) UpdateAboutSectionImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Image string `json:"image"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	img, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.New("invalid image"))
		return
	}

	err = controller.projects.UpdateAboutSectionImage(ctx, id, bytes.NewBuffer(img))
	if err != nil {
		controller.log.Error("could not create about section image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// GetAboutSection is an endpoint that returns about section of collection mint site.
func (controller *Projects) GetAboutSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	aboutSection, err := controller.projects.GetAboutSection(ctx, id)
	if err != nil {
		if collections.ErrNoAboutSection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not get about section", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(aboutSection); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// DeleteAboutSectionImage is an endpoint that deletes about section image.
func (controller *Projects) DeleteAboutSectionImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteAboutSectionImage(ctx, id)
	if err != nil {
		controller.log.Error("could not delete about section image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// AddMember is an endpoint that adds member of collection.
func (controller *Projects) AddMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.AddMember(ctx, id, req.Name, req.Role)
	if err != nil {
		controller.log.Error("could not create member", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateMemberImage is an endpoint that updates member image.
func (controller *Projects) UpdateMemberImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["memberId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Image string `json:"image"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	img, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.New("invalid image"))
		return
	}

	err = controller.projects.UpdateMemberImage(ctx, id, bytes.NewBuffer(img))
	if err != nil {
		if collections.ErrNoMember.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not create member image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// ListMembers is an endpoint that returns all members of collection site.
func (controller *Projects) ListMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	members, err := controller.projects.ListMembers(ctx, id)
	if err != nil {
		controller.log.Error("could not list members", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(members); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// UpdateMember is an endpoint that updates certain member.
func (controller *Projects) UpdateMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["memberId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.UpdateMember(ctx, collections.Member{
		ID:   id,
		Name: req.Name,
		Role: req.Role,
	})
	if err != nil {
		if collections.ErrNoMember.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not update member", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// DeleteMember is an endpoint that deletes certain member.
func (controller *Projects) DeleteMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["memberId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteMember(ctx, id)
	if err != nil {
		if collections.ErrNoMember.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not delete member", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// DeleteMemberImage is an endpoint that deletes member image.
func (controller *Projects) DeleteMemberImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	memberID, err := uuid.Parse(vars["memberId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteMemberImage(ctx, projectID, memberID)
	if err != nil {
		controller.log.Error("could not delete member image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// UpdateContentImage is an endpoint that updates content image.
func (controller *Projects) UpdateContentImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Image string `json:"image"`
		Name  string `json:"name"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	img, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.New("invalid image"))
		return
	}

	err = controller.projects.UpdateContentImage(ctx, id, req.Name, bytes.NewBuffer(img))
	if err != nil {
		controller.log.Error("could not create content image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// ListContentImages is an endpoint that returns names of content images.
func (controller *Projects) ListContentImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	names, err := controller.projects.ListContentFiles(ctx, id)
	if err != nil {
		if localfilestore.ErrNoFile.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not list content images", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(names); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// DeleteContentImage is an endpoint that deletes content image.
func (controller *Projects) DeleteContentImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	name := strings.TrimSpace(vars["name"])
	if name == "" {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.New("invalid image name"))
		return
	}

	err = controller.projects.DeleteContentImage(ctx, id, name)
	if err != nil {
		controller.log.Error("could not delete content image", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// AddMilestone is an endpoint that adds milestone.
func (controller *Projects) AddMilestone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"date"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.AddMilestone(ctx, id, req.Date, req.Title, req.Description)
	if err != nil {
		controller.log.Error("could not create milestone", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// ListMilestones is an endpoint that returns all collection milestones.
func (controller *Projects) ListMilestones(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	milestones, err := controller.projects.ListMilestones(ctx, id)
	if err != nil {
		controller.log.Error("could not list milestones", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(milestones); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// UpdateMilestone is an endpoint that updates certain milestone.
func (controller *Projects) UpdateMilestone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["milestoneId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"date"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.UpdateMilestone(ctx, collections.Milestone{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Date:        req.Date,
	})
	if err != nil {
		if collections.ErrNoMilestone.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not update milestone", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// DeleteMilestone is an endpoint that deletes milestone.
func (controller *Projects) DeleteMilestone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["milestoneId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteMilestone(ctx, id)
	if err != nil {
		if collections.ErrNoMilestone.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not delete milestone", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// AddSocialNetwork is an endpoint that adds collection social network.
func (controller *Projects) AddSocialNetwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.AddSocialNetwork(ctx, id, req.Title, req.Link)
	if err != nil {
		controller.log.Error("could not create social network", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// ListSocialNetworks is an endpoint that returns all collection social networks.
func (controller *Projects) ListSocialNetworks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	socialNetworks, err := controller.projects.ListSocialNetworks(ctx, id)
	if err != nil {
		controller.log.Error("could not list social networks", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(socialNetworks); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// UpdateSocialNetwork is an endpoint that updates social network.
func (controller *Projects) UpdateSocialNetwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["socialNetworkId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	type request struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.log.Error("failed to read json response", ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.UpdateSocialNetwork(ctx, collections.SocialNetwork{
		ID:    id,
		Title: req.Title,
		Link:  req.Link,
	})
	if err != nil {
		if collections.ErrNoSocialNetwork.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not update social network", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// DeleteSocialNetwork is an endpoint that deletes social network.
func (controller *Projects) DeleteSocialNetwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["socialNetworkId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteSocialNetwork(ctx, id)
	if err != nil {
		if collections.ErrNoSocialNetwork.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		}
		controller.log.Error("could not delete social network", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
		return
	}
}

// GetProjectPage renders the mint site preview for the collection owner.
func (controller *Projects) GetProjectPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	project, err := controller.projects.Get(ctx, id)
	if err != nil {
		controller.serveSiteError(w, err)
		return
	}

	site, err := controller.projects.GetSite(ctx, id)
	if err != nil {
		controller.serveSiteError(w, err)
		return
	}

	controller.renderSite(w, site, project)
}

// PublicSitePage renders the published mint site at /sites/{slug}; it does not require authentication.
func (controller *Projects) PublicSitePage(w http.ResponseWriter, r *http.Request) {
	site, project, err := controller.projects.GetPublishedSite(r.Context(), mux.Vars(r)["slug"])
	if err != nil {
		controller.serveSiteError(w, err)
		return
	}

	controller.renderSite(w, site, project)
}

// renderSite executes the mint site template with the collection contract and network for its scripts.
func (controller *Projects) renderSite(w http.ResponseWriter, site collections.Site, project collections.Collection) {
	if controller.template == nil {
		controller.log.Error("project template is not set", ErrProjects.New("project template is not set"))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.New("project template is not set"))
		return
	}

	page := mintSitePage{
		Site: site,
		Config: mintSiteConfig{
			ProjectID:       project.ID,
			ContractAddress: project.ContractAddress.Hex(),
			NFTCount:        project.NFTCount,
			MintPrice:       nfts.DefaultMintPrice,
			MaxPerMint:      nfts.DefaultMaxPerMint,
			Chain:           controller.chain,
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := controller.template.Execute(w, page); err != nil {
		controller.log.Error("could not execute get project template", ErrProjects.Wrap(err))
	}
}

// serveSiteError responds to a failed mint site request.
func (controller *Projects) serveSiteError(w http.ResponseWriter, err error) {
	switch {
	case collections.ErrNoCollection.Has(err), collections.ErrNoHeaderSection.Has(err), localfilestore.ErrNoFile.Has(err):
		controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
	default:
		controller.log.Error("could not get project site", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
	}
}

// PublishSite publishes the collection mint site at /sites/{slug}.
func (controller *Projects) PublishSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id, userID, ok := controller.projectAndUser(w, r)
	if !ok {
		return
	}

	var request struct {
		Slug string `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	project, err := controller.projects.PublishSite(ctx, id, userID, strings.TrimSpace(request.Slug))
	if err != nil {
		controller.serveWebsiteError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(project); err != nil {
		controller.log.Error("failed to write json response", ErrProjects.Wrap(err))
		return
	}
}

// UnpublishSite hides the collection mint site.
func (controller *Projects) UnpublishSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, userID, ok := controller.projectAndUser(w, r)
	if !ok {
		return
	}

	if err := controller.projects.UnpublishSite(r.Context(), id, userID); err != nil {
		controller.serveWebsiteError(w, err)
		return
	}
}

// projectAndUser parses the collection id from the path and the user id from the access token.
func (controller *Projects) projectAndUser(w http.ResponseWriter, r *http.Request) (id, userID uuid.UUID, ok bool) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return id, userID, false
	}

	claims, err := pkgauth.GetTokenClaims(r.Context())
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrProjects.Wrap(err))
		return id, userID, false
	}

	userID, err = uuid.Parse(claims["userId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return id, userID, false
	}

	return id, userID, true
}

// serveWebsiteError responds to a failed publish or unpublish request.
func (controller *Projects) serveWebsiteError(w http.ResponseWriter, err error) {
	switch {
	case collections.ErrInvalidSlug.Has(err), collections.ErrContractNotDeployed.Has(err), collections.ErrNoHeaderSection.Has(err):
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
	case collections.ErrCollectionOwnership.Has(err):
		controller.serveError(w, http.StatusForbidden, ErrProjects.Wrap(err))
	case collections.ErrNoCollection.Has(err):
		controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
	case collections.ErrSlugTaken.Has(err):
		controller.serveError(w, http.StatusConflict, ErrProjects.Wrap(err))
	default:
		controller.log.Error("could not change project website", ErrProjects.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
	}
}

// DeleteProjectSite is an endpoint that deletes project site.
func (controller *Projects) DeleteProjectSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrProjects.Wrap(err))
		return
	}

	err = controller.projects.DeleteSite(ctx, id)
	if err != nil {
		switch {
		case collections.ErrNoSite.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		case localfilestore.ErrNoFile.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrProjects.Wrap(err))
			return
		default:
			controller.log.Error("could not delete project site", ErrProjects.Wrap(err))
			controller.serveError(w, http.StatusInternalServerError, ErrProjects.Wrap(err))
			return
		}
	}
}

// serveError replies to request with specific code and error.
func (controller *Projects) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", ErrProjects.Wrap(err))
	}
}
