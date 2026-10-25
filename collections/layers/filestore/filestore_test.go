package filestore_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/collections/layers"
	"generator/collections/layers/filestore"
)

func TestFileStore(t *testing.T) {
	tmpDirectoryOutput, err := ioutil.TempDir("testdata", "output")
	require.NoError(t, err)

	tmpDirectoryInput, err := ioutil.TempDir("testdata", "input")
	require.NoError(t, err)

	defer func() {
		err = os.RemoveAll(tmpDirectoryInput)
		require.NoError(t, err)

		err = os.RemoveAll(tmpDirectoryOutput)
		require.NoError(t, err)
	}()

	store := filestore.NewStore(filestore.Config{
		LayersPath: tmpDirectoryInput,
		OutputPath: tmpDirectoryOutput,
	})

	id := uuid.New()
	ctx := context.Background()

	t.Run("save generated image", func(t *testing.T) {
		testImgPath := "testdata/testimg.png"

		img, err := os.Open(testImgPath)
		require.NoError(t, err)

		err = store.SaveImage(ctx, id, "1.png", img)
		require.NoError(t, err)
	})

	t.Run("create layer image", func(t *testing.T) {
		testImgPath := "testdata/testimg.png"

		img, err := os.Open(testImgPath)
		require.NoError(t, err)

		err = store.Create(ctx, id, "testlayer", "1.png", img)
		require.NoError(t, err)
	})

	t.Run("GetImageName", func(t *testing.T) {
		name, err := store.GetImageName(ctx, id, "testlayer", 1)
		require.NoError(t, err)

		require.Equal(t, "1.png", name)
	})

	t.Run("Negative GetImageName", func(t *testing.T) {
		_, err := store.GetImageName(ctx, id, "testlayer", 3)
		require.Error(t, err)
		assert.True(t, layers.ErrNoLayerFile.Has(err))
	})

	t.Run("get layer image", func(t *testing.T) {
		reader, err := store.GetLayerImage(ctx, id, "testlayer", 1)
		require.NoError(t, err)

		img, err := ioutil.ReadAll(reader)
		require.NoError(t, err)

		expectedImg, err := ioutil.ReadFile("testdata/testimg.png")
		require.NoError(t, err)

		compareImages(t, img, expectedImg)
	})

	t.Run("list layer images", func(t *testing.T) {
		layerImages, err := store.ListLayerImages(ctx, id, "testlayer")
		require.NoError(t, err)
		require.Equal(t, len(layerImages), 1)

		assert.Equal(t, layerImages[0], "1.png")
	})

	t.Run("update", func(t *testing.T) {
		err := store.UpdateName(ctx, "1.png", "2.png", id, "testlayer")
		require.NoError(t, err)

		name, err := store.GetImageName(ctx, id, "testlayer", 1)
		require.NoError(t, err)

		require.Equal(t, "2.png", name)
	})

	t.Run("delete", func(t *testing.T) {
		err := store.Delete(ctx, id, "testlayer", "2.png")
		require.NoError(t, err)
	})

	t.Run("delete folder with content", func(t *testing.T) {
		err = store.Create(ctx, id, "testlayer", "1.png", bytes.NewBuffer([]byte{}))
		require.NoError(t, err)

		err = store.DeleteFolder(ctx, id, "testlayer")
		require.NoError(t, err)
	})
}

func compareImages(t *testing.T, expectedImg, actualImg []byte) {
	require.Equal(t, len(expectedImg), len(actualImg))

	for i := 0; i < len(expectedImg); i++ {
		assert.Equal(t, expectedImg[i], actualImg[i])
	}
}
