package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/collections"
	"generator/internal/logger/zaplog"
	pkgauth "generator/pkg/auth"
)

// fakeCollections stores collections in memory.
type fakeCollections map[uuid.UUID]collections.Collection

func (fake fakeCollections) Get(ctx context.Context, id uuid.UUID) (collections.Collection, error) {
	project, ok := fake[id]
	if !ok {
		return project, collections.ErrNoCollection.New("")
	}
	return project, nil
}

func TestWithCollectionOwner(t *testing.T) {
	owner, stranger := uuid.New(), uuid.New()
	project := collections.Collection{ID: uuid.New(), UserID: owner}

	// the test token middleware stands for withAuth: the user id comes in the X-User header.
	withUser := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if user := r.Header.Get("X-User"); user != "" {
				r = r.WithContext(pkgauth.SetTokenClaims(r.Context(), pkgauth.TokenClaims{"userId": user}))
			}
			handler.ServeHTTP(w, r)
		})
	}
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	router := mux.NewRouter()
	collectionsRouter := router.PathPrefix("/collections").Subrouter()
	collectionsRouter.Use(withUser)
	collectionsRouter.Use(withCollectionOwner(zaplog.NewLog(), fakeCollections{project.ID: project}))
	collectionsRouter.Handle("", ok)
	collectionsRouter.Handle("/{id}", ok)
	collectionsRouter.PathPrefix("/{projectId}/mint-site").Subrouter().Handle("/header-section", ok)

	tests := []struct {
		name   string
		path   string
		user   uuid.UUID
		status int
	}{
		{"owner", "/collections/" + project.ID.String(), owner, http.StatusOK},
		{"owner, nested route", "/collections/" + project.ID.String() + "/mint-site/header-section", owner, http.StatusOK},
		{"stranger", "/collections/" + project.ID.String(), stranger, http.StatusForbidden},
		{"stranger, nested route", "/collections/" + project.ID.String() + "/mint-site/header-section", stranger, http.StatusForbidden},
		{"unknown collection", "/collections/" + uuid.New().String(), owner, http.StatusNotFound},
		{"invalid id", "/collections/1/mint-site/header-section", owner, http.StatusBadRequest},
		{"no collection in the path", "/collections", stranger, http.StatusOK},
		{"no token", "/collections/" + project.ID.String(), uuid.Nil, http.StatusUnauthorized},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			if test.user != uuid.Nil {
				request.Header.Set("X-User", test.user.String())
			}
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			require.Equal(t, test.status, recorder.Code, recorder.Body.String())
			if test.status != http.StatusOK {
				assert.Contains(t, recorder.Body.String(), `"error"`)
			}
		})
	}
}
