package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"generator/collections"
	"generator/internal/logger"
	pkgauth "generator/pkg/auth"
	"generator/pkg/pagination"
	"generator/whitelist"
)

// ErrWhitelist is an internal error type for whitelist controller.
var ErrWhitelist = errs.Class("whitelist controller error")

// Whitelist is a mvc controller that handles all whitelist related views.
type Whitelist struct {
	log       logger.Logger
	whitelist *whitelist.Service
	projects  *collections.Service
}

// NewWhitelist is a constructor for whitelist controller.
func NewWhitelist(log logger.Logger, whitelist *whitelist.Service, projects *collections.Service) *Whitelist {
	whitelistController := &Whitelist{
		log:       log,
		whitelist: whitelist,
		projects:  projects,
	}

	return whitelistController
}

// Add is an endpoint that adds address to whitelist.
func (controller *Whitelist) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	claims, err := pkgauth.GetTokenClaims(ctx)
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrWhitelist.Wrap(err))
		return
	}

	userID, err := uuid.Parse(claims["userId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	request := struct {
		ProjectID     uuid.UUID `json:"projectId"`
		Name          string    `json:"name"`
		WalletAddress string    `json:"walletAddress"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	if strings.TrimSpace(request.Name) == "" {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.New("name is empty"))
		return
	}

	if strings.TrimSpace(request.Name) == "" {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.New("name is empty"))
		return
	}

	if !common.IsHexAddress(request.WalletAddress) {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.New("invalid wallet address"))
		return
	}

	project, err := controller.projects.Get(ctx, request.ProjectID)
	if err != nil {
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrWhitelist.Wrap(err))
			return
		}
		controller.serveError(w, http.StatusInternalServerError, ErrWhitelist.Wrap(err))
		return
	}

	if project.UserID != userID {
		controller.serveError(w, http.StatusForbidden, ErrWhitelist.New("this project does not belong to the user"))
		return
	}

	item := whitelist.Item{
		Name:          request.Name,
		WalletAddress: common.HexToAddress(request.WalletAddress),
	}
	if err = controller.whitelist.Add(ctx, request.ProjectID, item); err != nil {
		controller.log.Error("could not add wallet address to whitelist", ErrWhitelist.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrWhitelist.Wrap(err))
		return
	}
}

// AddFromExcelFile is an endpoint that adds wallet addresses to whitelist from excel file.
func (controller *Whitelist) AddFromExcelFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	claims, err := pkgauth.GetTokenClaims(ctx)
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrWhitelist.Wrap(err))
		return
	}

	userID, err := uuid.Parse(claims["userId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	request := struct {
		File string `json:"file"`
	}{}
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	if request.File == "" {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.New("whitelist file not loaded"))
		return
	}

	file, err := base64.StdEncoding.DecodeString(request.File)
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	fileReader := bytes.NewReader(file)
	if err = controller.whitelist.AddFromExcelFile(ctx, projectID, userID, fileReader); err != nil {
		if whitelist.ErrInvalidWhitelistFile.Has(err) {
			controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
			return
		}
		if collections.ErrCollectionOwnership.Has(err) {
			controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
			return
		}
		controller.log.Error("could not add wallet address to whitelist from excel file", ErrWhitelist.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrWhitelist.Wrap(err))
		return
	}
}

// GetMerkleProof is an endpoint that returns merkle proof by wallet address.
func (controller *Whitelist) GetMerkleProof(w http.ResponseWriter, r *http.Request) {
	controller.serveMerkleProof(w, r, controller.whitelist.GetMerkleProof)
}

// GetPublicMerkleProof is an endpoint that returns Merkle proof for a buyer on the published mint site;
// it does not require authentication.
func (controller *Whitelist) GetPublicMerkleProof(w http.ResponseWriter, r *http.Request) {
	controller.serveMerkleProof(w, r, controller.whitelist.GetPublicMerkleProof)
}

// GetMerkleRoot is an endpoint that returns the Merkle root of the collection whitelist: the deployed
// contract must hold it (setMerkleRoot) for whitelisted buyers to mint; zero for an empty whitelist.
func (controller *Whitelist) GetMerkleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projectID, err := uuid.Parse(mux.Vars(r)["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	merkleRoot, err := controller.whitelist.GetMerkleRoot(r.Context(), projectID)
	if err != nil {
		controller.log.Error("could not get merkle root", ErrWhitelist.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrWhitelist.Wrap(err))
		return
	}

	response := struct {
		MerkleRoot string `json:"merkleRoot"`
	}{MerkleRoot: "0x" + hex.EncodeToString(merkleRoot)}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json response", ErrWhitelist.Wrap(err))
	}
}

// serveMerkleProof responds with the Merkle proof of the wallet address as an array of hex strings;
// the proof of the only address in the whitelist is empty.
func (controller *Whitelist) serveMerkleProof(w http.ResponseWriter, r *http.Request,
	getProof func(ctx context.Context, projectID uuid.UUID, walletAddress common.Address) ([][]byte, error)) {
	merkleProofStr := []string{}
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	if !common.IsHexAddress(vars["walletAddress"]) {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.New("invalid address"))
		return
	}
	walletAddress := common.HexToAddress(vars["walletAddress"])

	merkleProof, err := getProof(ctx, projectID, walletAddress)
	if err != nil {
		if whitelist.ErrNoItem.Has(err) || collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrWhitelist.Wrap(err))
			return
		}
		controller.log.Error("could not get merkle proof by wallet address", ErrWhitelist.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrWhitelist.Wrap(err))
		return
	}

	for _, proof := range merkleProof {
		merkleProofStr = append(merkleProofStr, "0x"+hex.EncodeToString(proof))
	}

	if err = json.NewEncoder(w).Encode(merkleProofStr); err != nil {
		controller.log.Error("failed to write json response", ErrWhitelist.Wrap(err))
		return
	}
}

// List is an endpoint that returns whitelist by project id.
func (controller *Whitelist) List(w http.ResponseWriter, r *http.Request) {
	var (
		cursor        pagination.Cursor
		err           error
		whitelistPage whitelist.Page
	)

	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	vars := mux.Vars(r)
	parametersQuery := r.URL.Query()
	limitParameter := parametersQuery.Get("limit")
	pageParameter := parametersQuery.Get("page")

	if limitParameter != "" {
		if cursor.Limit, err = strconv.Atoi(limitParameter); err != nil {
			controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
			return
		}
	}

	if pageParameter != "" {
		if cursor.Page, err = strconv.Atoi(pageParameter); err != nil {
			controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
			return
		}
	}

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	name := parametersQuery.Get(string(whitelist.FilterName))
	if name == "" {
		whitelistPage, err = controller.whitelist.ListByProjectIDWithPagination(ctx, projectID, cursor)
	} else {
		whitelistPage, err = controller.whitelist.ListByProjectIDAndItemName(ctx, projectID, name, cursor)
	}
	if err != nil {
		controller.log.Error("could not list whitelist", ErrWhitelist.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrWhitelist.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(whitelistPage); err != nil {
		controller.log.Error("failed to write json response", ErrWhitelist.Wrap(err))
		return
	}
}

// Delete is an endpoint that deletes wallet address from whitelist.
func (controller *Whitelist) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	vars := mux.Vars(r)

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrWhitelist.Wrap(err))
		return
	}

	walletAddress := common.HexToAddress(vars["walletAddress"])

	if err := controller.whitelist.Delete(ctx, projectID, walletAddress); err != nil {
		if whitelist.ErrNoItem.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrWhitelist.Wrap(err))
			return
		}
		controller.log.Error("could not delete wallet address from whitelist", ErrWhitelist.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrWhitelist.Wrap(err))
		return
	}
}

// serveError replies to request with specific code and error.
func (controller *Whitelist) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", ErrWhitelist.Wrap(err))
	}
}
