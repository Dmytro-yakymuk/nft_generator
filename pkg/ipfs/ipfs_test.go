package ipfs_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/pkg/ipfs"
)

func TestAddDirectory(t *testing.T) {
	ctx := context.Background()
	const directory = "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"

	received := make(map[string]string)
	node := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v0/add", r.URL.Path)
		for name, value := range map[string]string{"wrap-with-directory": "true", "cid-version": "1", "pin": "true", "quieter": "true"} {
			assert.Equal(t, value, r.URL.Query().Get(name), name)
		}

		reader, err := r.MultipartReader()
		require.NoError(t, err)
		for {
			part, err := reader.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			require.NoError(t, err)
			content, err := io.ReadAll(part)
			require.NoError(t, err)
			received[part.FileName()] = string(content)
		}
		_, _ = fmt.Fprintf(w, `{"Name":"","Hash":"%s","Size":"42"}`+"\n", directory)
	}))
	defer node.Close()

	file := func(name, content string) ipfs.File {
		return ipfs.File{Name: name, Content: func() ([]byte, error) { return []byte(content), nil }}
	}

	cid, err := ipfs.NewClient(node.URL+"/").AddDirectory(ctx, []ipfs.File{file("0.png", "image 0"), file("1", `{"name":"#1"}`)})
	require.NoError(t, err)
	assert.Equal(t, directory, cid)
	assert.Equal(t, map[string]string{"0.png": "image 0", "1": `{"name":"#1"}`}, received)

	t.Run("node error", func(t *testing.T) {
		failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(io.Discard, r.Body)
			http.Error(w, "repo is locked", http.StatusInternalServerError)
		}))
		defer failing.Close()

		_, err := ipfs.NewClient(failing.URL).AddDirectory(ctx, []ipfs.File{file("0.png", "image 0")})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "repo is locked")
	})

	t.Run("file can not be read", func(t *testing.T) {
		// the request body breaks off, so this node only drains it.
		draining := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(io.Discard, r.Body)
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer draining.Close()

		broken := ipfs.File{Name: "0.png", Content: func() ([]byte, error) { return nil, errors.New("no such file") }}
		_, err := ipfs.NewClient(draining.URL).AddDirectory(ctx, []ipfs.File{broken})
		require.Error(t, err)
	})

	t.Run("not configured", func(t *testing.T) {
		_, err := ipfs.NewClient("").AddDirectory(ctx, nil)
		require.Error(t, err)
	})
}

func TestExportCAR(t *testing.T) {
	ctx := context.Background()
	const directory = "bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"

	node := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v0/dag/export", r.URL.Path)
		if r.URL.Query().Get("arg") != directory {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"Message":"block was not found locally (offline)"}`))
			return
		}
		w.Header().Set("Trailer", "X-Stream-Error")
		_, _ = w.Write([]byte("car bytes"))
		if r.URL.Query().Get("broken") != "" {
			w.Header().Set("X-Stream-Error", "unexpected EOF")
		}
	}))
	defer node.Close()

	var car bytes.Buffer
	require.NoError(t, ipfs.NewClient(node.URL).ExportCAR(ctx, directory, &car))
	assert.Equal(t, "car bytes", car.String())

	err := ipfs.NewClient(node.URL).ExportCAR(ctx, "bafkreimissing", io.Discard)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found locally")

	assert.Error(t, ipfs.NewClient("").ExportCAR(ctx, directory, io.Discard))
}
