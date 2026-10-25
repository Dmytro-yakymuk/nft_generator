package layers_test

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"generator/collections/layers"
	"generator/collections/layers/filestore"
	"generator/collections/metadata"
	"generator/internal/filestore/localfilestore"
	"generator/internal/logger/zaplog"
)

// writeLayers makes 4 layers of 3 variants of 64x64: a background and semi-transparent discs in different
// places; the last layer is optional.
func writeLayers(t *testing.T, root string, projectID uuid.UUID) []layers.Layer {
	var projectLayers []layers.Layer
	for k := 0; k < 4; k++ {
		layer := layers.Layer{ID: uuid.New(), ProjectID: projectID, Name: "layer" + strconv.Itoa(k), Priority: k, Weights: map[string]int{}}
		dir := filepath.Join(root, projectID.String(), layer.Name)
		require.NoError(t, os.MkdirAll(dir, 0o755))
		for v := 0; v < 3; v++ {
			img := image.NewNRGBA(image.Rect(0, 0, 64, 64))
			for y := 0; y < 64; y++ {
				for x := 0; x < 64; x++ {
					dx, dy := x-12-k*12, y-12-v*16
					if k == 0 || dx*dx+dy*dy < 100 {
						img.Set(x, y, color.NRGBA{R: uint8(80 * v), G: uint8(60 * k), B: uint8(x * 4), A: uint8(120 + 40*v)})
					}
				}
			}
			name := "v" + strconv.Itoa(v+1) + ".png"
			file, err := os.Create(filepath.Join(dir, name))
			require.NoError(t, err)
			require.NoError(t, png.Encode(file, img))
			require.NoError(t, file.Close())
			layer.ImageNames = append(layer.ImageNames, name)
		}
		layer.LayerImagesCount = len(layer.ImageNames)
		if k == 3 {
			layer.NoneWeight = 100
		}
		projectLayers = append(projectLayers, layer)
	}
	return projectLayers
}

// TestGenerateImagesParallel checks that a pool of workers makes exactly the files of one worker.
func TestGenerateImagesParallel(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	projectID := uuid.New()
	projectLayers := writeLayers(t, root, projectID)

	// every combination of 3 * 3 * 3 * 4 = 108, in a scrambled order.
	var combinations []*big.Int
	for i := 0; i < 108; i++ {
		combinations = append(combinations, big.NewInt(int64((i*37)%108+1)))
	}

	generate := func(workers int) string {
		out := t.TempDir()
		service := layers.NewService(zaplog.NewLog(), layers.Config{Width: 64, Height: 64, Workers: workers}, nil,
			metadata.NewService(metadata.Config{OutputPath: "metadata"}, localfilestore.NewFileStore(localfilestore.Config{ImagesPath: out})),
			filestore.NewStore(filestore.Config{LayersPath: root, OutputPath: filepath.Join(out, "images")}),
			localfilestore.NewFileStore(localfilestore.Config{ImagesPath: out}),
		)
		require.NoError(t, service.GenerateImages(ctx, projectID, combinations, projectLayers))
		return out
	}

	sequential, parallel := generate(1), generate(4)
	for _, dir := range []string{"images", "metadata"} {
		files, err := os.ReadDir(filepath.Join(sequential, dir, projectID.String()))
		require.NoError(t, err)
		require.Len(t, files, len(combinations), dir)
		for _, file := range files {
			want, err := os.ReadFile(filepath.Join(sequential, dir, projectID.String(), file.Name()))
			require.NoError(t, err)
			got, err := os.ReadFile(filepath.Join(parallel, dir, projectID.String(), file.Name()))
			require.NoError(t, err)
			require.Equal(t, want, got, "%s/%s", dir, file.Name())
		}
	}
}

// TestPNGCompression checks that a faster compression level changes only the files, not the pixels.
func TestPNGCompression(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	projectID := uuid.New()
	projectLayers := writeLayers(t, root, projectID)
	combinations := []*big.Int{big.NewInt(1), big.NewInt(50), big.NewInt(107)}

	generate := func(compression string) string {
		out := t.TempDir()
		service := layers.NewService(zaplog.NewLog(), layers.Config{Width: 64, Height: 64, Workers: 2, PNGCompression: compression}, nil,
			metadata.NewService(metadata.Config{OutputPath: "metadata"}, localfilestore.NewFileStore(localfilestore.Config{ImagesPath: out})),
			filestore.NewStore(filestore.Config{LayersPath: root, OutputPath: filepath.Join(out, "images")}),
			localfilestore.NewFileStore(localfilestore.Config{ImagesPath: out}),
		)
		require.NoError(t, service.GenerateImages(ctx, projectID, combinations, projectLayers))
		return filepath.Join(out, "images", projectID.String())
	}

	decode := func(path string) *image.NRGBA {
		file, err := os.Open(path)
		require.NoError(t, err)
		defer func() { require.NoError(t, file.Close()) }()
		img, err := png.Decode(file)
		require.NoError(t, err)
		nrgba := image.NewNRGBA(img.Bounds())
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				nrgba.Set(x, y, img.At(x, y))
			}
		}
		return nrgba
	}

	standard, fast := generate("default"), generate("best-speed")
	for token := range combinations {
		name := strconv.Itoa(token) + ".png"
		require.Equal(t, decode(filepath.Join(standard, name)).Pix, decode(filepath.Join(fast, name)).Pix, name)
	}

	for _, compression := range []string{"", "default", "best-speed", "best-compression", "none"} {
		_, err := layers.Config{PNGCompression: compression}.PNGEncoder()
		require.NoError(t, err, compression)
	}
	_, err := layers.Config{PNGCompression: "fastest"}.PNGEncoder()
	require.Error(t, err)
}
