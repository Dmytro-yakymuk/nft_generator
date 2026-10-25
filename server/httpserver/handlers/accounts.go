package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/accounts"
	"generator/internal/logger"
	pkgauth "generator/pkg/auth"
)

// ErrUsers is an internal error type for users controller.
var ErrUsers = errs.Class("users controller error")

// Users is a mvc controller that handles all users related views.
type Users struct {
	log   logger.Logger
	users *accounts.Service
}

// NewUsers is a constructor for users controller.
func NewUsers(log logger.Logger, users *accounts.Service) *Users {
	usersController := &Users{
		log:   log,
		users: users,
	}

	return usersController
}

// Get is an endpoint that returns user.
func (controller *Users) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	claims, err := pkgauth.GetTokenClaims(ctx)
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrUsers.Wrap(err))
		return
	}

	userID, err := uuid.Parse(claims["userId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrUsers.Wrap(err))
		return
	}

	user, err := controller.users.Get(ctx, userID)
	if err != nil {
		if accounts.ErrNoAccount.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrUsers.Wrap(err))
			return
		}
		controller.log.Error("could not get user", ErrUsers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrUsers.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		controller.log.Error("failed to write json response", ErrUsers.Wrap(err))
		return
	}
}

// Update is an endpoint that updates user.
func (controller *Users) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	claims, err := pkgauth.GetTokenClaims(ctx)
	if err != nil {
		controller.serveError(w, http.StatusUnauthorized, ErrUsers.Wrap(err))
		return
	}

	userID, err := uuid.Parse(claims["userId"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrUsers.Wrap(err))
		return
	}

	var userSetting accounts.AccountSettings
	if err = json.NewDecoder(r.Body).Decode(&userSetting); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrUsers.Wrap(err))
		return
	}
	userSetting.ID = userID

	user, err := controller.users.Get(ctx, userID)
	if err != nil {
		if accounts.ErrNoAccount.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrUsers.Wrap(err))
			return
		}
		controller.log.Error("could not get user", ErrUsers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrUsers.Wrap(err))
		return
	}
	userSetting.Status = user.Status

	if err = controller.users.Update(ctx, userSetting); err != nil {
		if accounts.ErrNoAccount.Has(err) {
			controller.serveError(w, http.StatusNotFound, ErrUsers.Wrap(err))
			return
		}
		controller.log.Error("could not update user name", ErrUsers.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrUsers.Wrap(err))
		return
	}
}

// serveError replies to request with specific code and error.
func (controller *Users) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", ErrUsers.Wrap(err))
	}
}
