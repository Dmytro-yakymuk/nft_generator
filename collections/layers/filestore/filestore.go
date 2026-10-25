package filestore

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/collections/layers"
)

// ensures that Store implements layers.Store.
var _ layers.Store = (*Store)(nil)

// ErrFileStore indicates that there was an error in the file store.
var ErrFileStore = errs.Class("FileStore repository error")

// Config defines configurable variables of Store.
type Config struct {
	LayersPath string `json:"layersPath"`
	OutputPath string `json:"outputPath"`
}

// Store describes layer storage in file system.
type Store struct {
	config Config
}

// NewStore is constructor for file system store.
func NewStore(cfg Config) *Store {
	return &Store{
		config: cfg,
	}
}

// MakeDirectory creates directory in the file system by specific path.
func (store *Store) MakeDirectory(ctx context.Context, path string) error {
	outputPath := filepath.Join(store.config.LayersPath, path)

	return ErrFileStore.Wrap(os.MkdirAll(outputPath, os.ModePerm))
}

// SaveImage saves generated image.
func (store *Store) SaveImage(ctx context.Context, projectID uuid.UUID, imageName string, reader io.Reader) (err error) {
	err = os.MkdirAll(filepath.Join(store.config.OutputPath, projectID.String()), os.ModePerm)
	if err != nil {
		return ErrFileStore.Wrap(err)
	}

	file, err := os.Create(filepath.Join(store.config.OutputPath, projectID.String(), imageName))
	if err != nil {
		return ErrFileStore.Wrap(err)
	}

	defer func() {
		err = ErrFileStore.Wrap(file.Close())
	}()

	_, err = io.Copy(file, reader)

	return ErrFileStore.Wrap(err)
}

// Create writes layer in the file system.
func (store *Store) Create(ctx context.Context, projectID uuid.UUID, layerName, layerImgName string, reader io.Reader) (err error) {
	outputPath := filepath.Join(store.config.LayersPath, projectID.String(), layerName)

	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return ErrFileStore.Wrap(err)
	}

	file, err := os.Create(filepath.Join(outputPath, layerImgName))
	if err != nil {
		return ErrFileStore.Wrap(err)
	}
	defer func() {
		err = ErrFileStore.Wrap(file.Close())
	}()

	_, err = io.Copy(file, reader)

	return ErrFileStore.Wrap(err)
}

// GetImageName returns image name by sequence number of image.
func (store *Store) GetImageName(ctx context.Context, projectID uuid.UUID, layerName string, imgSequenceNumber int) (string, error) {
	path := filepath.Join(store.config.LayersPath, projectID.String(), layerName)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", layers.ErrNoLayerFile.Wrap(err)
		}
		return "", ErrFileStore.Wrap(err)
	}

	if len(files) < imgSequenceNumber-1 {
		return "", layers.ErrNoLayerFile.New("")
	}

	return files[imgSequenceNumber-1].Name(), ErrFileStore.Wrap(err)
}

// GetLayerImage returns layer image by sequence number of image.
func (store *Store) GetLayerImage(ctx context.Context, projectID uuid.UUID, layerName string, imgSequenceNumber int) (io.Reader, error) {
	path := filepath.Join(store.config.LayersPath, projectID.String(), layerName)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, layers.ErrNoLayerFile.Wrap(err)
		}
		return nil, ErrFileStore.Wrap(err)
	}

	if len(files) < imgSequenceNumber-1 {
		return nil, layers.ErrNoLayerFile.New("")
	}

	output := filepath.Join(path, files[imgSequenceNumber-1].Name())
	file, err := os.Open(output)

	return file, ErrFileStore.Wrap(err)
}

// ListLayerImages returns all names of layer images.
func (store *Store) ListLayerImages(ctx context.Context, projectID uuid.UUID, layerName string) ([]string, error) {
	path := filepath.Join(store.config.LayersPath, projectID.String(), layerName)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, layers.ErrNoLayerFile.Wrap(err)
		}
		return nil, ErrFileStore.Wrap(err)
	}

	var names []string
	for _, file := range files {
		names = append(names, file.Name())
	}

	return names, ErrFileStore.Wrap(err)
}

// UpdateName updates name of layer.
func (store *Store) UpdateName(ctx context.Context, oldName, newName string, projectID uuid.UUID, layerName string) error {
	pathWithOldName := filepath.Join(store.config.LayersPath, projectID.String(), layerName, oldName)
	pathWithNewName := filepath.Join(store.config.LayersPath, projectID.String(), layerName, newName)

	return ErrFileStore.Wrap(os.Rename(pathWithOldName, pathWithNewName))
}

// Delete deletes layer image from file system.
func (store *Store) Delete(ctx context.Context, projectID uuid.UUID, layerName string, layerImageName string) error {
	outputPath := filepath.Join(store.config.LayersPath, projectID.String(), layerName, layerImageName)

	return ErrFileStore.Wrap(os.Remove(outputPath))
}

// DeleteFolder deletes layer folder with content inside.
func (store *Store) DeleteFolder(ctx context.Context, projectID uuid.UUID, layerName string) error {
	path := filepath.Join(store.config.LayersPath, projectID.String(), layerName)

	return ErrFileStore.Wrap(os.RemoveAll(path))
}
