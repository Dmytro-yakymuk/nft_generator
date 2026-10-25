package images

import (
	"context"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/internal/filestore"
	"generator/internal/filestore/localfilestore"
	"generator/pkg/storage"
)

// ErrImages indicates that there was an error in the service.
var ErrImages = errs.Class("images service error")

// Service is handling images related logic.
//
// architecture: Service
type Service struct {
	config    Config
	filestore filestore.FileStore
}

// NewService is a constructor for images service.
func NewService(config Config, filestore filestore.FileStore) *Service {
	return &Service{
		config:    config,
		filestore: filestore,
	}
}

// List returns the names of the generated images in the order of token ids. The names are numbers (token ids,
// or combination numbers in collections generated before), so they are sorted by value rather than as text,
// where 10.png would come before 2.png.
func (service *Service) List(ctx context.Context, projectID uuid.UUID) ([]string, error) {
	path := filepath.Join(service.config.ImagesPath, projectID.String())
	names, err := service.filestore.List(ctx, path)
	sort.SliceStable(names, func(a, b int) bool { return lessByNumber(names[a], names[b]) })
	return names, ErrImages.Wrap(err)
}

// lessByNumber compares file names whose base names are numbers without leading zeros by value: a shorter
// number is smaller, numbers of equal length compare as text. Other names compare as text.
func lessByNumber(a, b string) bool {
	x, y := strings.TrimSuffix(a, filepath.Ext(a)), strings.TrimSuffix(b, filepath.Ext(b))
	if len(x) != len(y) && isNumber(x) && isNumber(y) {
		return len(x) < len(y)
	}
	return x < y
}

func isNumber(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Get returns the content of the generated image.
func (service *Service) Get(ctx context.Context, projectID uuid.UUID, name string) ([]byte, error) {
	path := filepath.Join(service.config.ImagesPath, projectID.String())
	content, err := service.filestore.Get(ctx, path, name)
	return content, ErrImages.Wrap(err)
}

// GeneratedImagesCount returns count of generated images.
func (service *Service) GeneratedImagesCount(ctx context.Context, projectID uuid.UUID) (int, error) {
	path := filepath.Join(service.config.ImagesPath, projectID.String())
	count, err := service.filestore.Count(ctx, path)
	return count, ErrImages.Wrap(err)
}

// Upload uploads image to file store.
func (service *Service) Upload(ctx context.Context, projectID uuid.UUID, oldName string) error {
	relatedPath := filepath.Join(service.config.ImagesPath, projectID.String())

	client, err := storage.NewClient(service.config.FileStorage)
	if err != nil {
		return ErrImages.Wrap(err)
	}

	file, err := service.filestore.Get(ctx, relatedPath, oldName)
	if err != nil {
		return ErrImages.Wrap(err)
	}

	fullPath := filepath.Join(projectID.String(), oldName)
	err = client.Upload(ctx, service.config.Bucket, fullPath, file)

	return ErrImages.Wrap(err)
}

// DeleteGeneratedImages deletes all denerated images for project from directory.
func (service *Service) DeleteGeneratedImages(ctx context.Context, projectID uuid.UUID) error {
	relatedPath := filepath.Join(service.config.ImagesPath, projectID.String())

	if err := service.filestore.DeleteFolder(ctx, relatedPath); err != nil {
		if !localfilestore.ErrNoFile.Has(err) {
			return ErrImages.Wrap(err)
		}
	}

	return nil
}

// DeleteProjectImages deletes all images related to project.
func (service *Service) DeleteProjectImages(ctx context.Context, projectID uuid.UUID) error {
	return ErrImages.Wrap(service.filestore.DeleteFolder(ctx, projectID.String()))
}
