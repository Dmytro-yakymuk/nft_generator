package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"generator/collections"
	"generator/internal/logger"
	"generator/nfts"
	pkgauth "generator/pkg/auth"
	"generator/whitelist"
)

// ErrContracts is an internal error type for contracts controller.
var ErrContracts = errs.Class("contracts controller error")

// Contracts is a mvc controller that handles all contracts related views.
type Contracts struct {
	log  logger.Logger
	nfts *nfts.Service
}

// NewContracts is a constructor for contracts controller.
func NewContracts(log logger.Logger, nfts *nfts.Service) *Contracts {
	contractsController := &Contracts{
		log:  log,
		nfts: nfts,
	}

	return contractsController
}

// GetDeployContractTransactionData is an endpoint that returns contract data.
func (controller *Contracts) GetDeployContractTransactionData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	vars := mux.Vars(r)

	claimsMap, err := pkgauth.GetTokenClaims(ctx)
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrContracts.Wrap(err))
		return
	}

	userID, err := uuid.Parse(claimsMap["userId"])
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrContracts.Wrap(err))
		return
	}

	projectID, err := uuid.Parse(vars["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrContracts.Wrap(err))
		return
	}

	// the sale terms come from the contract page; missing ones take the defaults.
	query := r.URL.Query()
	terms := nfts.Terms{Name: query.Get("name"), Price: query.Get("price")}
	if maxPerMint := query.Get("maxPerMint"); maxPerMint != "" {
		if terms.MaxPerMint, err = strconv.Atoi(maxPerMint); err != nil || terms.MaxPerMint < 1 {
			controller.serveError(w, http.StatusBadRequest, ErrContracts.New("maxPerMint must be a positive integer"))
			return
		}
	}

	data, err := controller.nfts.GetDeployContractTransactionData(ctx, projectID, userID, terms)
	if err != nil {
		if nfts.ErrInvalidTerms.Has(err) {
			controller.serveError(w, http.StatusBadRequest, ErrContracts.Wrap(err))
			return
		}
		if collections.ErrNoCollection.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrContracts.Wrap(err))
			return
		}
		if whitelist.ErrNoItem.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrContracts.Wrap(err))
			return
		}
		controller.log.Error("could not get data", ErrContracts.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrContracts.Wrap(err))
		return
	}

	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)

	if err = json.NewEncoder(w).Encode(string(dst)); err != nil {
		controller.log.Error("failed to write json response", ErrContracts.Wrap(err))
		return
	}
}

// serveError replies to request with specific code and error.
func (controller *Contracts) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", ErrContracts.Wrap(err))
	}
}
