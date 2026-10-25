package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleMintSiteHosts(t *testing.T) {
	serve := func(domain, host, path string) string {
		router := mux.NewRouter()
		handleMintSiteHosts(router, domain, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("site " + mux.Vars(r)["slug"]))
		})
		// the web app takes every other request, as in the server.
		router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("app"))
		})

		request := httptest.NewRequest(http.MethodGet, path, nil)
		request.Host = host
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		return recorder.Body.String()
	}

	assert.Equal(t, "site cosmic-critters", serve("localhost", "cosmic-critters.localhost:9012", "/"))
	assert.Equal(t, "site cosmic-critters", serve("nft.example.com", "cosmic-critters.nft.example.com", "/"))
	// the platform itself, other paths of a site host and hosts of other domains go to the web app.
	assert.Equal(t, "app", serve("localhost", "localhost:9012", "/"))
	assert.Equal(t, "app", serve("localhost", "cosmic-critters.localhost:9012", "/static/js/index.js"))
	assert.Equal(t, "app", serve("nft.example.com", "cosmic-critters.example.com", "/"))
	assert.Equal(t, "app", serve("localhost", "a.b.localhost", "/"))
	assert.Equal(t, "app", serve("", "cosmic-critters.localhost:9012", "/"))
}
