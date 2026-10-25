package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"generator/collections"
	"generator/dashboard"
	"generator/internal/logger"
	pkgauth "generator/pkg/auth"
	"generator/pkg/chain"
)

// ErrDashboard is an internal error type for dashboard controller.
var ErrDashboard = errs.Class("dashboard controller error")

// Dashboard is a mvc controller that handles all dashboard related views.
type Dashboard struct {
	log       logger.Logger
	dashboard *dashboard.Service
}

// NewDashboard is a constructor for dashboard controller.
func NewDashboard(log logger.Logger, dashboard *dashboard.Service) *Dashboard {
	dashboardController := &Dashboard{
		log:       log,
		dashboard: dashboard,
	}

	return dashboardController
}

// Events is an endpoint that returns a page of the collection contract events, newest first.
func (controller *Dashboard) Events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projectID, userID, ok := controller.projectAndUser(w, r)
	if !ok {
		return
	}

	page := 1
	if pageArg := r.URL.Query().Get("page"); pageArg != "" {
		var err error
		if page, err = strconv.Atoi(pageArg); err != nil || page < 1 {
			controller.serveError(w, http.StatusBadRequest, ErrDashboard.New("invalid page"))
			return
		}
	}

	events, err := controller.dashboard.Events(r.Context(), projectID, userID, page)
	if err != nil {
		controller.serveServiceError(w, "could not get events", err)
		return
	}

	if err = json.NewEncoder(w).Encode(events); err != nil {
		controller.log.Error("failed to write json response", ErrDashboard.Wrap(err))
		return
	}
}

// CollectionStats is an endpoint that returns the number of tokens, owners and sales of the collection.
func (controller *Dashboard) CollectionStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projectID, userID, ok := controller.projectAndUser(w, r)
	if !ok {
		return
	}

	stats, err := controller.dashboard.CollectionStats(r.Context(), projectID, userID)
	if err != nil {
		controller.serveServiceError(w, "could not get collection statistic", err)
		return
	}

	if err = json.NewEncoder(w).Encode(stats); err != nil {
		controller.log.Error("failed to write json response", ErrDashboard.Wrap(err))
		return
	}
}

// EventData is an endpoint that returns paid events and their prices for the period
// (the last `period` days, or between the unix times `from` and `to`).
func (controller *Dashboard) EventData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projectID, userID, ok := controller.projectAndUser(w, r)
	if !ok {
		return
	}

	var err error
	period := dashboard.Period
	from := int64(dashboard.Period)
	to := int64(dashboard.Period)
	query := r.URL.Query()
	periodArg := query.Get("period")
	if periodArg != "" {
		period, err = strconv.Atoi(periodArg)
		if err != nil || period < 0 {
			controller.serveError(w, http.StatusBadRequest, ErrDashboard.New("invalid period"))
			return
		}
	}
	fromDate := query.Get("from")
	if fromDate != "" {
		from, err = strconv.ParseInt(fromDate, 10, 64)
		if err != nil {
			controller.serveError(w, http.StatusBadRequest, ErrDashboard.New("invalid from date"))
			return
		}
	}
	toDate := query.Get("to")
	if toDate != "" {
		to, err = strconv.ParseInt(toDate, 10, 64)
		if err != nil {
			controller.serveError(w, http.StatusBadRequest, ErrDashboard.New("invalid to date"))
			return
		}
	}

	if periodArg != "" && fromDate != "" && toDate != "" {
		controller.serveError(w, http.StatusBadRequest, ErrDashboard.New("invalid input value"))
		return
	}

	eventDataAndAveragePrices, err := controller.dashboard.EventDataAndAveragePrices(r.Context(), projectID, userID, period, from, to)
	if err != nil {
		controller.serveServiceError(w, "could not get event data", err)
		return
	}

	if err = json.NewEncoder(w).Encode(eventDataAndAveragePrices); err != nil {
		controller.log.Error("failed to write json response", ErrDashboard.Wrap(err))
		return
	}
}

// projectAndUser parses the collection id from the path and the user id from the access token.
func (controller *Dashboard) projectAndUser(w http.ResponseWriter, r *http.Request) (projectID, userID uuid.UUID, ok bool) {
	projectID, err := uuid.Parse(mux.Vars(r)["projectId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrDashboard.Wrap(err))
		return projectID, userID, false
	}

	claims, err := pkgauth.GetTokenClaims(r.Context())
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrDashboard.Wrap(err))
		return projectID, userID, false
	}

	userID, err = uuid.Parse(claims["userId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrDashboard.Wrap(err))
		return projectID, userID, false
	}

	return projectID, userID, true
}

// serveServiceError replies with the status that matches the dashboard service error.
func (controller *Dashboard) serveServiceError(w http.ResponseWriter, message string, err error) {
	switch {
	case collections.ErrNoCollection.Has(err):
		controller.serveError(w, http.StatusNotFound, ErrDashboard.Wrap(err))
	case collections.ErrCollectionOwnership.Has(err):
		controller.serveError(w, http.StatusForbidden, ErrDashboard.Wrap(err))
	case collections.ErrContractNotDeployed.Has(err):
		controller.serveError(w, http.StatusConflict, ErrDashboard.Wrap(err))
	case chain.ErrChain.Has(err):
		// the node is unavailable or rejected the request.
		controller.log.Error(message, ErrDashboard.Wrap(err))
		controller.serveError(w, http.StatusBadGateway, ErrDashboard.Wrap(err))
	default:
		controller.log.Error(message, ErrDashboard.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrDashboard.Wrap(err))
	}
}

// serveError replies to request with specific code and error.
func (controller *Dashboard) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", ErrDashboard.Wrap(err))
	}
}
