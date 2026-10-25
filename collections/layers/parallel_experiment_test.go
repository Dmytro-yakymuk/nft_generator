package layers

import (
	"context"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"generator/collections/metadata"
	"generator/internal/filestore/localfilestore"
	"generator/internal/logger/zaplog"
)

// artStore reads the layer images of the art set from data/demo-art/<layer>/ and saves the generated images
// to a directory, like the file system store of the server (which this package cannot import in tests).
type artStore struct {
	art, out string
	layers   map[string][]string
}

func (store *artStore) GetLayerImage(ctx context.Context, projectID uuid.UUID, layerName string, number int) (io.Reader, error) {
	return os.Open(filepath.Join(store.art, layerName, store.layers[layerName][number-1]))
}

func (store *artStore) SaveImage(ctx context.Context, projectID uuid.UUID, name string, reader io.Reader) error {
	file, err := os.Create(filepath.Join(store.out, name))
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	return errs.Combine(err, file.Close())
}

func (store *artStore) MakeDirectory(context.Context, string) error { return nil }
func (store *artStore) Create(context.Context, uuid.UUID, string, string, io.Reader) error {
	return nil
}
func (store *artStore) GetImageName(context.Context, uuid.UUID, string, int) (string, error) {
	return "", nil
}
func (store *artStore) ListLayerImages(ctx context.Context, projectID uuid.UUID, layerName string) ([]string, error) {
	return store.layers[layerName], nil
}
func (store *artStore) UpdateName(context.Context, string, string, uuid.UUID, string) error {
	return nil
}
func (store *artStore) Delete(context.Context, uuid.UUID, string, string) error { return nil }
func (store *artStore) DeleteFolder(context.Context, uuid.UUID, string) error   { return nil }

// TestExperimentParallel measures the generation of 500 images of the 26-layer art set by 1 to 12 workers
// (the processor has 8 cores and 12 threads), with the layer images decoded once per generation.
func TestExperimentParallel(t *testing.T) {
	dir := experimentDir(t)
	projectLayers := artLayers(t)
	art := filepath.Join("..", "..", "data", "demo-art")
	projectID := uuid.New()
	names := make(map[string][]string, len(projectLayers))
	for k := range projectLayers {
		projectLayers[k].ProjectID = projectID
		names[projectLayers[k].Name] = projectLayers[k].ImageNames
	}

	const amount = 500
	combinations, err := chooseCombinations(rand.New(rand.NewSource(1)), amount, projectLayers)
	require.NoError(t, err)

	var rows [][]string
	for _, workers := range []int{1, 2, 4, 6, 8, 10, 12} {
		for repeat := 1; repeat <= experimentRepeats; repeat++ {
			out := t.TempDir()
			service := NewService(zaplog.NewLog(), Config{Width: 500, Height: 500, Workers: workers}, nil,
				metadata.NewService(metadata.Config{OutputPath: "metadata"}, localfilestore.NewFileStore(localfilestore.Config{ImagesPath: out})),
				&artStore{art: art, out: out, layers: names}, nil)

			runtime.GC()
			start := time.Now()
			require.NoError(t, service.GenerateImages(context.Background(), projectID, combinations, projectLayers))
			elapsed := time.Since(start).Seconds()

			rows = append(rows, []string{
				strconv.Itoa(workers), strconv.Itoa(runtime.NumCPU()), strconv.Itoa(amount), strconv.Itoa(repeat),
				strconv.FormatFloat(elapsed, 'f', 3, 64), strconv.FormatFloat(elapsed/amount*1000, 'f', 2, 64),
			})
			t.Logf("workers=%d #%d: %.2f s, %.1f ms per image", workers, repeat, elapsed, elapsed/amount*1000)
		}
	}

	writeCSV(t, filepath.Join(dir, "parallel.csv"), []string{"workers", "cpus", "amount", "repeat", "seconds", "ms_per_image"}, rows)
}
