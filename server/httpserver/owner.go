package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"generator/collections"
	"generator/internal/logger"
	pkgauth "generator/pkg/auth"
)

// collectionOwners returns collections with their owners; *collections.Service implements it.
type collectionOwners interface {
	Get(ctx context.Context, id uuid.UUID) (collections.Collection, error)
}

// withCollectionOwner lets a request reach the route only if the collection named in the path
// ({projectId} or {id}) belongs to the signed-in user; requests without a collection in the path pass as is.
// It must run after withAuth, which puts the access token claims into the request context.
func withCollectionOwner(log logger.Logger, projects collectionOwners) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			idArg, ok := vars["projectId"]
			if !ok {
				idArg, ok = vars["id"]
			}
			if !ok {
				handler.ServeHTTP(w, r)
				return
			}

			id, err := uuid.Parse(idArg)
			if err != nil {
				serveJSONError(w, http.StatusBadRequest, "invalid collection id")
				return
			}

			claims, err := pkgauth.GetTokenClaims(r.Context())
			if err != nil {
				serveJSONError(w, http.StatusUnauthorized, "authorization required")
				return
			}
			userID, err := uuid.Parse(claims["userId"])
			if err != nil {
				serveJSONError(w, http.StatusUnauthorized, "authorization required")
				return
			}

			project, err := projects.Get(r.Context(), id)
			if err != nil {
				if collections.ErrNoCollection.Has(err) {
					serveJSONError(w, http.StatusNotFound, "collection does not exist")
					return
				}
				log.Error("could not check the collection owner", Error.Wrap(err))
				serveJSONError(w, http.StatusInternalServerError, "could not check the collection owner")
				return
			}
			if project.UserID != userID {
				serveJSONError(w, http.StatusForbidden, "this collection does not belong to user")
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}

// serveJSONError replies with the status and {"error": message}, as the handlers do.
func serveJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
