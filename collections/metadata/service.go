package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/internal/filestore"
	"generator/internal/filestore/localfilestore"
	"generator/internal/nft"
	"generator/pkg/storage"
)

// ErrMetadata indicates that there was an error in the service.
var ErrMetadata = errs.Class("metadata service error")

// Service is handling metadata related logic.
//
// architecture: Service
type Service struct {
	config Config
	store  filestore.FileStore
}

// NewService is a constructor for metadata service.
func NewService(config Config, store filestore.FileStore) *Service {
	return &Service{
		config: config,
		store:  store,
	}
}

// GenerateMetadata generates metadata of nft token in open sea formate for the image with the given number
// (the token id; collections generated before were named by the combination number); dna is the number of
// the layer combination of the image and the attributes keep the layer order.
func (service *Service) GenerateMetadata(ctx context.Context, projectID uuid.UUID, number *big.Int, dna string, attributes []nft.Attribute) error {
	metadata := nft.Metadata{
		Description: service.config.Metadata.Description,
		ExternalURL: service.config.Metadata.ExternalURL,
		Image:       fmt.Sprintf(service.config.Metadata.Image, projectID.String(), number),
		Name:        service.config.Metadata.Name,
		Attributes:  attributes,
		DNA:         dna,
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return ErrMetadata.Wrap(err)
	}
	reader := bytes.NewReader(data)
	relatedPath := filepath.Join(service.config.OutputPath, projectID.String(), number.String())

	return ErrMetadata.Wrap(service.store.Store(ctx, relatedPath, reader))
}

// Upload uploads metadate file to file store.
func (service *Service) Upload(ctx context.Context, projectID uuid.UUID, oldName string, newName string) error {
	result := strings.Split(oldName, ".")
	name := result[0]
	relatedPath := filepath.Join(service.config.OutputPath, projectID.String())

	client, err := storage.NewClient(service.config.FileStorage)
	if err != nil {
		return ErrMetadata.Wrap(err)
	}

	file, err := service.store.Get(ctx, relatedPath, name)
	if err != nil {
		return ErrMetadata.Wrap(err)
	}

	fullPath := filepath.Join(projectID.String(), newName)
	err = client.Upload(ctx, service.config.Bucket, fullPath, file)

	return ErrMetadata.Wrap(err)
}

// Get returns stored metadata of the generated image.
func (service *Service) Get(ctx context.Context, projectID uuid.UUID, imageName string) (nft.Metadata, error) {
	name := strings.Split(imageName, ".")[0]
	relatedPath := filepath.Join(service.config.OutputPath, projectID.String())

	data, err := service.store.Get(ctx, relatedPath, name)
	if err != nil {
		return nft.Metadata{}, ErrMetadata.Wrap(err)
	}

	var metadata nft.Metadata
	if err = json.Unmarshal(data, &metadata); err != nil {
		return nft.Metadata{}, ErrMetadata.Wrap(err)
	}

	// The link is built from the current config, so metadata generated before
	// the public address was set still points to the reachable image.
	if number, ok := new(big.Int).SetString(name, 10); ok {
		metadata.Image = fmt.Sprintf(service.config.Metadata.Image, projectID.String(), number)
	}

	return metadata, nil
}

// Delete deletes metadata files for project from file store.
func (service *Service) Delete(ctx context.Context, projectID uuid.UUID) error {
	relatedPath := filepath.Join(service.config.OutputPath, projectID.String())
	if err := service.store.DeleteFolder(ctx, relatedPath); err != nil {
		if !localfilestore.ErrNoFile.Has(err) {
			return ErrMetadata.Wrap(err)
		}
	}

	return nil
}
