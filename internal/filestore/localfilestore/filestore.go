package localfilestore

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zeebo/errs"

	"generator/internal/filestore"
)

// ensures that FileStore implements filestore.FileStore.
var _ filestore.FileStore = (*LocalFileStore)(nil)

// ErrNoFile indicated that file does not exist.
var ErrNoFile = errs.Class("file does not exist")

// Config defines configurable values for file store.
type Config struct {
	ImagesPath string `json:"imagesPath"`
}

// LocalFileStore implements FileStore in file system.
type LocalFileStore struct {
	config Config
}

// NewFileStore is constructor for FileStore.
func NewFileStore(cfg Config) *LocalFileStore {
	return &LocalFileStore{
		config: cfg,
	}
}

// Store stores file in the local file system store.
func (localFileStore *LocalFileStore) Store(ctx context.Context, relatedPath string, reader io.Reader) error {
	relatedPath, fileName := filepath.Split(relatedPath)
	if fileName == "" {
		return errs.New("file name is empty")
	}

	path := filepath.Join(localFileStore.config.ImagesPath, relatedPath)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)

	return errs.Combine(err, file.Close())
}

// Get returns file by name from local file system store.
func (localFileStore *LocalFileStore) Get(ctx context.Context, relatedPath string, name string) ([]byte, error) {
	fullPath := filepath.Join(localFileStore.config.ImagesPath, relatedPath, name)
	file, err := ioutil.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoFile.Wrap(err)
		}
		return nil, err
	}

	return file, nil
}

// List returns all names of files from local file system store.
func (localFileStore *LocalFileStore) List(ctx context.Context, relatedPath string) ([]string, error) {
	path := filepath.Join(localFileStore.config.ImagesPath, relatedPath)
	filesInfo, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoFile.Wrap(err)
		}
		return nil, err
	}

	var names []string

	for _, fileInfo := range filesInfo {
		names = append(names, fileInfo.Name())
	}

	return names, nil
}

// Count returns count of files from local file system store.
func (localFileStore *LocalFileStore) Count(ctx context.Context, relatedPath string) (int, error) {
	path := filepath.Join(localFileStore.config.ImagesPath, relatedPath)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrNoFile.Wrap(err)
		}
		return 0, err
	}

	return len(files), nil
}

// Delete deletes file from local file system store.
func (localFileStore *LocalFileStore) Delete(ctx context.Context, relatedPath string) error {
	path := filepath.Join(localFileStore.config.ImagesPath, relatedPath)

	return os.Remove(path)
}

// DeleteFolder deletes folder from local file system store.
func (localFileStore *LocalFileStore) DeleteFolder(ctx context.Context, relatedPath string) error {
	path := filepath.Join(localFileStore.config.ImagesPath, relatedPath)

	if err := os.RemoveAll(path); err != nil {
		if os.IsNotExist(err) {
			return ErrNoFile.Wrap(err)
		}
		return err
	}

	return nil
}
