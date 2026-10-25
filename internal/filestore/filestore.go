package filestore

import (
	"context"
	"io"
)

// FileStore provides access to file storage.
type FileStore interface {
	// Store stores file in the file store.
	Store(ctx context.Context, relatedPath string, reader io.Reader) error
	// Get returns file by name from local file system store.
	Get(ctx context.Context, relatedPath string, name string) ([]byte, error)
	// List returns all names of file from directory.
	List(ctx context.Context, path string) ([]string, error)
	// Count returns count of files from local file system store.
	Count(ctx context.Context, relatedPath string) (int, error)
	// Delete deletes specific file from the file store.
	Delete(ctx context.Context, path string) error
	// DeleteFolder deletes folder from local file system store.
	DeleteFolder(ctx context.Context, relatedPath string) error
}
