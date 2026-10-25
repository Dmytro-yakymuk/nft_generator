package layers_test

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator"
	"generator/accounts"
	"generator/collections"
	"generator/collections/layers"
	"generator/collections/layers/filestore"
	"generator/collections/metadata"
	"generator/internal/filestore/localfilestore"
	"generator/internal/logger/zaplog"
	"generator/storage/testdb"
)

func TestLayersDB(t *testing.T) {
	user1 := accounts.Account{
		ID:           uuid.New(),
		Email:        "testUser1@gmail.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusCreated,
		CreatedAt:    time.Now().UTC(),
	}

	project1 := collections.Collection{
		ID:     uuid.New(),
		Name:   "project1",
		UserID: user1.ID,
	}

	layer1 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "face",
		Priority:  1,
	}

	layer2 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "eyes",
		Priority:  2,
	}

	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		repositoryUsers := db.Accounts()
		repositoryProjects := db.Collections()
		repositoryLayers := db.Layers()

		t.Run("Create", func(t *testing.T) {
			err := repositoryUsers.Create(ctx, user1)
			require.NoError(t, err)

			err = repositoryProjects.Create(ctx, project1)
			require.NoError(t, err)

			err = repositoryLayers.Create(ctx, layer1)
			require.NoError(t, err)

			err = repositoryLayers.Create(ctx, layer2)
			require.NoError(t, err)
		})

		t.Run("Get", func(t *testing.T) {
			layer, err := repositoryLayers.Get(ctx, layer1.ID)
			require.NoError(t, err)
			compareLayers(t, layer, layer1)
		})

		t.Run("List", func(t *testing.T) {
			layersFromDB, err := repositoryLayers.List(ctx, project1.ID)
			require.NoError(t, err)
			require.Equal(t, len(layersFromDB), 2)
			compareLayers(t, layersFromDB[0], layer1)
			compareLayers(t, layersFromDB[1], layer2)
		})

		t.Run("Update priority", func(t *testing.T) {
			layerToUpdate1 := layers.Layer{
				ID:       layer1.ID,
				Priority: 2,
			}

			layerToUpdate2 := layers.Layer{
				ID:       layer2.ID,
				Priority: 1,
			}

			err := repositoryLayers.UpdatePriorities(ctx, []layers.Layer{layerToUpdate1, layerToUpdate2})
			require.NoError(t, err)

			layer, err := repositoryLayers.Get(ctx, layer1.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 2)

			layer, err = repositoryLayers.Get(ctx, layer2.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 1)
		})

		t.Run("Weights", func(t *testing.T) {
			weights, err := repositoryLayers.ListWeights(ctx, project1.ID)
			require.NoError(t, err)
			require.Empty(t, weights)

			err = repositoryLayers.UpdateWeights(ctx, layer1.ID, 25, map[string]int{"round.png": 70, "square.png": 5})
			require.NoError(t, err)
			err = repositoryLayers.UpdateWeights(ctx, layer2.ID, 0, map[string]int{"happy.png": 3})
			require.NoError(t, err)

			layer, err := repositoryLayers.Get(ctx, layer1.ID)
			require.NoError(t, err)
			assert.Equal(t, 25, layer.NoneWeight)

			// the set replaces the former one.
			err = repositoryLayers.UpdateWeights(ctx, layer1.ID, 0, map[string]int{"round.png": 40})
			require.NoError(t, err)
			err = repositoryLayers.RenameWeight(ctx, layer1.ID, "round.png", "oval.png")
			require.NoError(t, err)
			err = repositoryLayers.DeleteWeight(ctx, layer2.ID, "happy.png")
			require.NoError(t, err)

			weights, err = repositoryLayers.ListWeights(ctx, project1.ID)
			require.NoError(t, err)
			assert.Equal(t, map[uuid.UUID]map[string]int{layer1.ID: {"oval.png": 40}}, weights)

			layersFromDB, err := repositoryLayers.List(ctx, project1.ID)
			require.NoError(t, err)
			for _, layer := range layersFromDB {
				assert.Equal(t, 0, layer.NoneWeight)
			}

			err = repositoryLayers.UpdateWeights(ctx, uuid.New(), 0, nil)
			require.True(t, layers.ErrNoLayer.Has(err))
		})

		t.Run("Delete", func(t *testing.T) {
			err := repositoryLayers.Delete(ctx, layer1.ID)
			require.NoError(t, err)

			// the weights go with the layer.
			weights, err := repositoryLayers.ListWeights(ctx, project1.ID)
			require.NoError(t, err)
			assert.Empty(t, weights)
		})
	})
}

func compareLayers(t *testing.T, layerFromDB, testLayer layers.Layer) {
	assert.Equal(t, layerFromDB.ID, testLayer.ID)
	assert.Equal(t, layerFromDB.ProjectID, testLayer.ProjectID)
	assert.Equal(t, layerFromDB.Name, testLayer.Name)
	assert.Equal(t, layerFromDB.Priority, testLayer.Priority)
}

func TestLayersService(t *testing.T) {
	user1 := accounts.Account{
		ID:           uuid.New(),
		Email:        "testUser1@gmail.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusCreated,
		CreatedAt:    time.Now().UTC(),
	}

	projectID, err := uuid.Parse("9751deb8-d0c9-493c-a500-b3030993cbfd")
	require.NoError(t, err)

	project1 := collections.Collection{
		ID:     projectID,
		Name:   "project1",
		UserID: user1.ID,
	}

	layer1 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "face",
		Priority:  1,
	}

	layer2 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "eye_brows",
		Priority:  2,
	}

	layer3 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "hair",
		Priority:  3,
	}

	layer4 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "lips",
		Priority:  4,
	}

	layer5 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "nose",
		Priority:  5,
	}

	layer6 := layers.Layer{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "testlayer",
		Priority:  6,
	}

	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		repositoryUsers := db.Accounts()
		repositoryProjects := db.Collections()
		repositoryLayers := db.Layers()

		tmpDirectory, err := ioutil.TempDir("testdata", "output")
		require.NoError(t, err)

		defer func() {
			err = os.RemoveAll(tmpDirectory)
			require.NoError(t, err)
		}()

		metadataService := metadata.NewService(
			metadata.Config{OutputPath: filepath.Base(tmpDirectory)},
			localfilestore.NewFileStore(localfilestore.Config{
				ImagesPath: "testdata",
			}),
		)

		layersService := layers.NewService(
			zaplog.NewLog(),
			layers.Config{},
			repositoryLayers,
			metadataService,
			filestore.NewStore(filestore.Config{
				LayersPath: "testdata",
				OutputPath: tmpDirectory,
			}),
			localfilestore.NewFileStore(localfilestore.Config{
				ImagesPath: "testdata",
			}),
		)

		t.Run("Create", func(t *testing.T) {
			err = repositoryUsers.Create(ctx, user1)
			require.NoError(t, err)

			err = repositoryProjects.Create(ctx, project1)
			require.NoError(t, err)

			err = repositoryLayers.Create(ctx, layer1)
			require.NoError(t, err)

			err = repositoryLayers.Create(ctx, layer2)
			require.NoError(t, err)

			err = repositoryLayers.Create(ctx, layer3)
			require.NoError(t, err)

			err = repositoryLayers.Create(ctx, layer4)
			require.NoError(t, err)

			err = repositoryLayers.Create(ctx, layer5)
			require.NoError(t, err)
		})

		t.Run("generate avatar", func(t *testing.T) {
			projectLayers, err := layersService.List(ctx, project1.ID)
			require.NoError(t, err)

			err = layersService.GenerateImages(ctx, project1.ID, []*big.Int{big.NewInt(1)}, projectLayers)
			require.NoError(t, err)

			generatedImg, err := ioutil.ReadFile(filepath.Join(tmpDirectory, project1.ID.String(), "0.png"))
			require.NoError(t, err)

			expectedImg, err := ioutil.ReadFile("testdata/9751deb8-d0c9-493c-a500-b3030993cbfd/result.png")
			require.NoError(t, err)

			compareImages(t, expectedImg, generatedImg)
		})

		t.Run("update priority raising", func(t *testing.T) {
			err = layersService.UpdatePriority(ctx, layer3.ID, 1)
			require.NoError(t, err)

			layer, err := repositoryLayers.Get(ctx, layer3.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 1)

			layer, err = repositoryLayers.Get(ctx, layer1.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 2)

			layer, err = repositoryLayers.Get(ctx, layer2.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 3)
		})

		t.Run("update priority downgrading", func(t *testing.T) {
			err = layersService.UpdatePriority(ctx, layer3.ID, 3)
			require.NoError(t, err)

			layer, err := repositoryLayers.Get(ctx, layer1.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 1)

			layer, err = repositoryLayers.Get(ctx, layer2.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 2)

			layer, err = repositoryLayers.Get(ctx, layer3.ID)
			require.NoError(t, err)
			assert.Equal(t, layer.Priority, 3)
		})

		t.Run("create layer image", func(t *testing.T) {
			err = repositoryLayers.Create(ctx, layer6)
			require.NoError(t, err)

			testImgPath := "testdata/testlayer/testimg.png"

			dat, err := os.Open(testImgPath)
			require.NoError(t, err)

			err = layersService.CreateLayerImage(ctx, layer6.ID, "testimg.png", dat)
			require.NoError(t, err)
		})

		t.Run("update name of layer image", func(t *testing.T) {
			err = layersService.UpdateNameOfLayerImage(ctx, layer6.ID, "1.png", "testimg.png")
			require.NoError(t, err)
		})

		t.Run("delete layer image", func(t *testing.T) {
			err = layersService.DeleteLayerImage(ctx, layer6.ID, "1.png")
			require.NoError(t, err)
		})
	})
}

func compareImages(t *testing.T, expectedImg, actualImg []byte) {
	require.Equal(t, len(expectedImg), len(actualImg))

	for i := 0; i < len(expectedImg); i++ {
		assert.Equal(t, expectedImg[i], actualImg[i])
	}
}
