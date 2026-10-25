package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zeebo/errs"

	"generator/accounts"
	authpkg "generator/accounts/auth"
	"generator/internal/logger"
	pkgauth "generator/pkg/auth"
)

// AuthError is a internal error for auth controller.
var AuthError = errs.Class("auth controller error")

// AuthTemplates holds all auth related templates.
type AuthTemplates struct {
	Login    *template.Template
	Register *template.Template
}

// Auth is an authentication controller that exposes users authentication functionality.
type Auth struct {
	log      logger.Logger
	userAuth *authpkg.Service
	cookie   *pkgauth.SessionCookie

	templates *AuthTemplates
}

// NewAuth returns new instance of Auth.
func NewAuth(log logger.Logger, userAuth *authpkg.Service, authCookie *pkgauth.SessionCookie, templates *AuthTemplates) *Auth {
	return &Auth{
		log:       log,
		userAuth:  userAuth,
		cookie:    authCookie,
		templates: templates,
	}
}

// Nonce is an endpoint to send nonce to metamask for login.
func (auth *Auth) Nonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	query := r.URL.Query()

	address := query.Get("address")
	if address == "" {
		auth.serveError(w, http.StatusBadRequest, AuthError.New("address is empty"))
		return
	}

	if !common.IsHexAddress(address) {
		auth.serveError(w, http.StatusBadRequest, AuthError.New("address is wrong"))
		return
	}

	nonce, err := auth.userAuth.Nonce(ctx, common.HexToAddress(address))
	if err != nil {
		switch {
		case accounts.ErrNoAccount.Has(err):
			auth.serveError(w, http.StatusNotFound, AuthError.Wrap(err))
		case authpkg.ErrUnauthenticated.Has(err):
			auth.serveError(w, http.StatusUnauthorized, AuthError.Wrap(err))
		default:
			auth.serveError(w, http.StatusInternalServerError, AuthError.Wrap(err))
		}
		return
	}

	if err = json.NewEncoder(w).Encode(nonce); err != nil {
		auth.log.Error("failed to write json response", AuthError.Wrap(err))
		return
	}
}

// MetamaskLogin is an endpoint to authorize user from metamask and set auth cookie in browser.
func (auth *Auth) MetamaskLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	type MetamaskFields struct {
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
	}

	var request MetamaskFields
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		auth.serveError(w, http.StatusBadRequest, AuthError.Wrap(err))
		return
	}

	if request.Signature == "" || request.Nonce == "" {
		auth.serveError(w, http.StatusBadRequest, AuthError.New("did not fill in all the fields"))
		return
	}

	signature, err := hexutil.Decode(request.Signature)
	if err != nil {
		auth.serveError(w, http.StatusBadRequest, AuthError.Wrap(err))
		return
	}

	authToken, err := auth.userAuth.LoginWithMetamask(ctx, request.Nonce, signature)
	if err != nil {
		switch {
		case authpkg.ErrUnauthenticated.Has(err):
			auth.serveError(w, http.StatusUnauthorized, AuthError.Wrap(err))
		default:
			auth.serveError(w, http.StatusInternalServerError, AuthError.Wrap(err))
		}

		return
	}

	auth.cookie.SetTokenCookie(w, r, authToken)
}

// MetamaskRegister is an endpoint to register user.
func (auth *Auth) MetamaskRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var sig string
	if err := json.NewDecoder(r.Body).Decode(&sig); err != nil {
		auth.serveError(w, http.StatusBadRequest, AuthError.Wrap(err))
		return
	}

	if sig == "" {
		auth.serveError(w, http.StatusBadRequest, AuthError.New("signature is empty"))
		return
	}

	signature, err := hexutil.Decode(sig)
	if err != nil {
		auth.serveError(w, http.StatusBadRequest, AuthError.Wrap(err))
		return
	}

	err = auth.userAuth.RegisterWithMetamask(ctx, signature)
	if err != nil {
		switch {
		case authpkg.ErrAccountExists.Has(err):
			auth.serveError(w, http.StatusConflict, AuthError.Wrap(err))
		case authpkg.ErrUnauthenticated.Has(err):
			auth.serveError(w, http.StatusUnauthorized, AuthError.Wrap(err))
		default:
			auth.serveError(w, http.StatusInternalServerError, AuthError.Wrap(err))
		}

		return
	}
}

// Logout is an endpoint to log out: it removes the auth cookie from the browser.
// It works without a valid token, so an expired session can be closed as well.
func (auth *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	auth.cookie.RemoveTokenCookie(w, r)
}

// serveError replies to request with specific code and error.
func (auth *Auth) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		auth.log.Error("failed to write json error response", AuthError.Wrap(err))
	}
}
