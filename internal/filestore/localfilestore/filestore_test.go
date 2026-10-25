package localfilestore_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/internal/filestore/localfilestore"
)

func TestFileStore(t *testing.T) {
	ctx := context.Background()

	tmpDir, err := ioutil.TempDir("testdata", "output")
	require.NoError(t, err)
	defer func() {
		err = os.RemoveAll(tmpDir)
		require.NoError(t, err)
	}()

	fileName := "testimg.png"
	relatedPath := filepath.Join("test", fileName)

	store := localfilestore.NewFileStore(localfilestore.Config{
		ImagesPath: tmpDir,
	})

	t.Run("store file", func(t *testing.T) {
		file, err := os.Open("testdata/testimg.png")
		require.NoError(t, err)

		err = store.Store(ctx, relatedPath, file)
		require.NoError(t, err)
	})

	t.Run("update file", func(t *testing.T) {
		file, err := os.Open("testdata/testimg.png")
		require.NoError(t, err)

		err = store.Store(ctx, relatedPath, file)
		require.NoError(t, err)
	})

	t.Run("get file", func(t *testing.T) {
		file, err := store.Get(ctx, "test", fileName)
		require.NoError(t, err)

		fullPath := filepath.Join(tmpDir, relatedPath)
		originFile, err := ioutil.ReadFile(fullPath)
		require.NoError(t, err)

		assert.Equal(t, originFile, file)
	})

	t.Run("list files names", func(t *testing.T) {
		names, err := store.List(ctx, "test")
		require.NoError(t, err)
		require.Equal(t, 1, len(names))

		assert.Equal(t, fileName, names[0])
	})

	t.Run("count", func(t *testing.T) {
		count, err := store.Count(ctx, "test")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("delete", func(t *testing.T) {
		err = store.Delete(ctx, relatedPath)
		require.NoError(t, err)
	})

	t.Run("list files names no files", func(t *testing.T) {
		names, err := store.List(ctx, "test")
		require.NoError(t, err)
		require.Equal(t, 0, len(names))
	})
}
