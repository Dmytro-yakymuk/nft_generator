package filebase_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"generator/pkg/filebase"
)

const cid = "bafybeidsypevgmrwdf6zeamm6t3yktmwveiua76whgkexpkv3zrqey7vpa"

func TestPinCAR(t *testing.T) {
	car := filepath.Join(t.TempDir(), "collection.car")
	content := []byte("car archive bytes")
	require.NoError(t, os.WriteFile(car, content, 0o600))
	digest := sha256.Sum256(content)

	// pinned is the root CID the fake bucket reports back.
	pinned := cid
	bucket := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/cosmic-critters/"+cid+".car", r.URL.Path)
		require.Equal(t, content, body)
		require.Equal(t, "car", r.Header.Get("X-Amz-Meta-Import"))
		require.Equal(t, hex.EncodeToString(digest[:]), r.Header.Get("X-Amz-Content-Sha256"))
		authorization := r.Header.Get("Authorization")
		require.True(t, strings.HasPrefix(authorization, "AWS4-HMAC-SHA256 Credential=KEY/"), authorization)
		require.Contains(t, authorization, "x-amz-meta-import")
		w.Header().Set("X-Amz-Meta-Cid", pinned)
	}))
	defer bucket.Close()

	client := filebase.NewClient(filebase.Config{Endpoint: bucket.URL, Bucket: "cosmic-critters", AccessKey: "KEY", SecretKey: "SECRET"})
	pin := func() error {
		file, err := os.Open(car)
		require.NoError(t, err)
		defer func() { require.NoError(t, file.Close()) }()
		return client.PinCAR(context.Background(), cid, file)
	}

	require.NoError(t, pin())

	pinned = "bafkreiother"
	require.Error(t, pin())
}

func TestPinCARRejected(t *testing.T) {
	bucket := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("<Error><Code>AccessDenied</Code></Error>"))
	}))
	defer bucket.Close()

	file, err := os.CreateTemp(t.TempDir(), "car")
	require.NoError(t, err)
	defer func() { require.NoError(t, file.Close()) }()

	err = filebase.NewClient(filebase.Config{Endpoint: bucket.URL, Bucket: "b", AccessKey: "k", SecretKey: "s"}).
		PinCAR(context.Background(), cid, file)
	require.Error(t, err)
	require.Contains(t, err.Error(), "AccessDenied")
}

func TestEnabled(t *testing.T) {
	require.False(t, filebase.Config{}.Enabled())
	require.False(t, filebase.Config{Bucket: "b", AccessKey: "k"}.Enabled())
	require.True(t, filebase.Config{Bucket: "b", AccessKey: "k", SecretKey: "s"}.Enabled())
}
