package httpserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrecompressed(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "dist"), 0o755))
	for name, content := range map[string]string{
		"dist/main.js": "plain", "dist/main.js.br": "brotli", "dist/main.js.gz": "gzip",
		"dist/logo.png": "png", "index.css": "css", "index.css.gz": "css gzip",
	} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
	}
	handler := newPrecompressed(dir)

	get := func(url, acceptEncoding string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodGet, url, nil)
		if acceptEncoding != "" {
			request.Header.Set("Accept-Encoding", acceptEncoding)
		}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		return recorder
	}

	for _, test := range []struct{ url, accept, body, encoding string }{
		{"/dist/main.js", "gzip, deflate, br, zstd", "brotli", "br"},
		{"/dist/main.js", "gzip", "gzip", "gzip"},
		{"/dist/main.js", "br;q=0, gzip;q=0.5", "gzip", "gzip"},
		{"/dist/main.js", "", "plain", ""},
		{"/dist/main.js", "identity", "plain", ""},
		{"/dist/logo.png", "br, gzip", "png", ""},
		{"/index.css", "br, gzip", "css gzip", "gzip"},
	} {
		response := get(test.url, test.accept)
		require.Equal(t, http.StatusOK, response.Code, "%+v", test)
		require.Equal(t, test.body, response.Body.String(), "%+v", test)
		require.Equal(t, test.encoding, response.Header().Get("Content-Encoding"), "%+v", test)
		require.Equal(t, "Accept-Encoding", response.Header().Get("Vary"), "%+v", test)
	}

	// the compressed copy keeps the type of the original file; hashed build files are kept for a year.
	response := get("/dist/main.js", "br")
	require.Contains(t, response.Header().Get("Content-Type"), "javascript")
	require.Contains(t, response.Header().Get("Cache-Control"), "immutable")
	require.Empty(t, get("/index.css", "").Header().Get("Cache-Control"))

	require.Equal(t, http.StatusNotFound, get("/missing.js", "br").Code)
}
