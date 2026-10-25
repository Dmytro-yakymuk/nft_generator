package layers

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"math/big"
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/collections/metadata"
	"generator/internal/filestore"
	"generator/internal/logger"
	"generator/internal/nft"
	"generator/pkg/imageproc"
)

// ErrLayers indicates that there was an error in the service.
var ErrLayers = errs.Class("layers service error")

// ErrCountLayers indicates that amount of images is greater than quantity of combination layers.
var ErrCountLayers = errs.Class("amount of images is greater than quantity of combination layers")

// Service is handling layers related logic.
//
// architecture: Service
type Service struct {
	log    logger.Logger
	config Config

	imageStore Store
	store      filestore.FileStore
	layers     DB
	metadata   *metadata.Service

	wg   sync.WaitGroup
	done chan struct{}

	mutex    sync.Mutex
	abortMap map[uuid.UUID]chan struct{}

	encoder *png.Encoder
}

// NewService is a constructor for layers service.
func NewService(log logger.Logger, config Config, layers DB, metadata *metadata.Service, imageStore Store, store filestore.FileStore) *Service {
	encoder, err := config.PNGEncoder()
	if err != nil {
		log.Error("generated images keep the default PNG compression", err)
		encoder = &png.Encoder{}
	}

	return &Service{
		log:        log,
		config:     config,
		layers:     layers,
		metadata:   metadata,
		imageStore: imageStore,
		store:      store,
		done:       make(chan struct{}),
		abortMap:   make(map[uuid.UUID]chan struct{}),
		encoder:    encoder,
	}
}

// Create creates new layer.
func (service *Service) Create(ctx context.Context, projectID uuid.UUID, name string, priority int) error {
	layer := Layer{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      name,
		Priority:  priority,
	}

	err := service.layers.Create(ctx, layer)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	path := filepath.Join(projectID.String(), name)

	return ErrLayers.Wrap(service.imageStore.MakeDirectory(ctx, path))
}

// List returns all project layers with their images and weights.
func (service *Service) List(ctx context.Context, projectID uuid.UUID) ([]Layer, error) {
	allLayers, err := service.layers.List(ctx, projectID)
	if err != nil {
		return nil, ErrLayers.Wrap(err)
	}

	weights, err := service.layers.ListWeights(ctx, projectID)
	if err != nil {
		return nil, ErrLayers.Wrap(err)
	}

	for i := range allLayers {
		layer := &allLayers[i]
		layer.ImageNames, err = service.imageStore.ListLayerImages(ctx, layer.ProjectID, layer.Name)
		if err != nil {
			return allLayers, ErrLayers.Wrap(err)
		}
		layer.LayerImagesCount = len(layer.ImageNames)

		// the client gets the weight of every image, the default one included.
		layer.Weights = weights[layer.ID]
		all := make(map[string]int, len(layer.ImageNames))
		for _, name := range layer.ImageNames {
			all[name] = layer.Weight(name)
		}
		layer.Weights = all
	}

	return allLayers, nil
}

// UpdateWeights sets the weight of the layer being absent (zero makes it mandatory) and the weights of its
// images; an image left out gets DefaultWeight.
func (service *Service) UpdateWeights(ctx context.Context, layerID uuid.UUID, noneWeight int, weights map[string]int) error {
	layer, err := service.layers.Get(ctx, layerID)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	if noneWeight < 0 || noneWeight > MaxWeight {
		return ErrInvalidWeight.New("weight of none %d is not from 0 to %d", noneWeight, MaxWeight)
	}

	names, err := service.imageStore.ListLayerImages(ctx, layer.ProjectID, layer.Name)
	if err != nil {
		return ErrLayers.Wrap(err)
	}
	images := make(map[string]bool, len(names))
	for _, name := range names {
		images[name] = true
	}
	for name, weight := range weights {
		if !images[name] {
			return ErrInvalidWeight.New("layer %q has no image %q", layer.Name, name)
		}
		if weight < 1 || weight > MaxWeight {
			return ErrInvalidWeight.New("weight %d of %q is not from 1 to %d", weight, name, MaxWeight)
		}
	}

	return ErrLayers.Wrap(service.layers.UpdateWeights(ctx, layerID, noneWeight, weights))
}

// UpdatePriority updates priority of layer.
func (service *Service) UpdatePriority(ctx context.Context, id uuid.UUID, priority int) error {
	layer, err := service.layers.Get(ctx, id)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	var index int

	switch {
	case priority > layer.Priority:
		index = -1
	case priority < layer.Priority:
		index = 1
	}

	allLayers, err := service.List(ctx, layer.ProjectID)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	var layersToUpdate []Layer

	for _, val := range allLayers {
		switch index {
		case 1:
			if val.Priority < layer.Priority && val.Priority >= priority {
				val.Priority += index
				layersToUpdate = append(layersToUpdate, val)
			}
		case -1:
			if val.Priority > layer.Priority && val.Priority <= priority {
				val.Priority += index
				layersToUpdate = append(layersToUpdate, val)
			}
		}
	}

	layer.Priority = priority
	layersToUpdate = append(layersToUpdate, layer)

	return ErrLayers.Wrap(service.layers.UpdatePriorities(ctx, layersToUpdate))
}

// ListLayerImageNames returns all names of layer images.
func (service *Service) ListLayerImageNames(ctx context.Context, projectID, layerID uuid.UUID) ([]string, error) {
	layer, err := service.layers.Get(ctx, layerID)
	if err != nil {
		return nil, ErrLayers.Wrap(err)
	}

	names, err := service.imageStore.ListLayerImages(ctx, projectID, layer.Name)
	return names, ErrLayers.Wrap(err)
}

// Delete deletes project layer and deletes layer photos.
func (service *Service) Delete(ctx context.Context, id uuid.UUID) error {
	layer, err := service.layers.Get(ctx, id)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	err = service.imageStore.DeleteFolder(ctx, layer.ProjectID, layer.Name)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	return ErrLayers.Wrap(service.layers.Delete(ctx, id))
}

// Generate chooses the combinations and starts images generation.
func (service *Service) Generate(ctx context.Context, projectID uuid.UUID, amount int) error {
	projectLayers, err := service.List(ctx, projectID)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	bases := layerBases(projectLayers)
	if amount > 0 && big.NewInt(int64(amount)).Cmp(combinationCount(bases)) > 0 {
		return ErrCountLayers.New("")
	}

	for _, layer := range projectLayers {
		for i := 1; i <= len(layer.ImageNames); i++ {
			reader, err := service.imageStore.GetLayerImage(ctx, layer.ProjectID, layer.Name, i)
			if err != nil {
				return ErrLayers.Wrap(err)
			}

			config, _, err := image.DecodeConfig(reader)
			closeReader(reader)
			if err != nil {
				return ErrLayers.Wrap(err)
			}

			if service.config.Width != config.Width || service.config.Height != config.Height {
				return ErrLayers.Wrap(ErrInvalidDimension)
			}
		}
	}

	// the combinations are chosen before the generation starts, so the client learns at once
	// if the weights leave too few distinct combinations.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	combinations, err := chooseCombinations(rnd, amount, projectLayers)
	if err != nil {
		return err
	}

	service.mutex.Lock()
	_, ok := service.abortMap[projectID]
	if ok {
		service.mutex.Unlock()
		return ErrGenerationProcess.New("generation process already exists")
	}
	abort := make(chan struct{})
	service.abortMap[projectID] = abort
	service.mutex.Unlock()

	service.wg.Add(1)
	go func() {
		defer service.wg.Done()
		// the project is unlocked however the generation ends: done, aborted or failed.
		defer func() {
			service.mutex.Lock()
			if service.abortMap[projectID] == abort {
				delete(service.abortMap, projectID)
			}
			service.mutex.Unlock()
		}()

		if err := service.generateImages(ctx, projectID, combinations, projectLayers, abort); err != nil {
			service.log.Error("could not generate image: ", ErrLayers.Wrap(err))
		}
	}()

	return nil
}

// chooseCombinations returns amount distinct combination numbers: by Floyd's algorithm if all variants of
// every layer are equally likely, otherwise by drawing every layer's variant by its weight.
func chooseCombinations(rnd *rand.Rand, amount int, projectLayers []Layer) ([]*big.Int, error) {
	weights := layerWeights(projectLayers)
	if isUniform(weights) {
		return sampleCombinations(rnd, amount, combinationCount(layerBases(projectLayers))), nil
	}

	combinations, draws, ok := sampleWeighted(rnd, amount, weights, maxDrawsPerImage*amount+1000)
	if !ok {
		return nil, ErrRareCombinations.New("%d of %d distinct combinations after %d draws: lower the collection size or make the weights more even",
			len(combinations), amount, draws)
	}
	return combinations, nil
}

// Close closes done channel.
func (service *Service) Close() {
	close(service.done)
	service.wg.Wait()
}

// GenerateImages generates the images of the given combinations, the i-th being the image of token i.
func (service *Service) GenerateImages(ctx context.Context, projectID uuid.UUID, combinations []*big.Int, projectLayers []Layer) error {
	service.mutex.Lock()
	abort := service.abortMap[projectID]
	service.mutex.Unlock()

	return service.generateImages(ctx, projectID, combinations, projectLayers, abort)
}

// generateImages generates the images by a pool of workers until all are done, abort is closed (a nil abort
// never fires) or an image fails; the images that are being made then are finished.
func (service *Service) generateImages(ctx context.Context, projectID uuid.UUID, combinations []*big.Int, projectLayers []Layer, abort <-chan struct{}) error {
	workers := service.config.Workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	if workers > len(combinations) {
		workers = len(combinations)
	}

	images := newLayerCache(service.imageStore, projectLayers)
	bases := layerBases(projectLayers)
	jobs := make(chan int)
	stop := make(chan struct{})
	var (
		wg       sync.WaitGroup
		failOnce sync.Once
		failure  error
	)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for token := range jobs {
				if err := service.generateImage(ctx, projectID, token, combinations[token], projectLayers, bases, images); err != nil {
					failOnce.Do(func() {
						failure = err
						close(stop)
					})
					return
				}
			}
		}()
	}

feed:
	for token := range combinations {
		select {
		case <-service.done:
			break feed
		case <-abort:
			break feed
		case <-stop:
			break feed
		case jobs <- token:
		}
	}
	close(jobs)
	wg.Wait()

	return failure
}

// generateImage composes the image of the token from the variants of the combination and saves it with its
// metadata. The metadata is saved first: the progress counts the images, so every counted image has it.
func (service *Service) generateImage(ctx context.Context, projectID uuid.UUID, token int, combination *big.Int, projectLayers []Layer, bases []int, images *layerCache) error {
	indices := layerIndices(combination, bases)
	parts := make([]image.Image, 0, len(projectLayers))
	canvas := image.Rect(0, 0, service.config.Width, service.config.Height)
	for k, layer := range projectLayers {
		// the index after the images is the "none" variant of an optional layer.
		if indices[k] > len(layer.ImageNames) {
			continue
		}

		part, full, err := images.get(ctx, k, indices[k])
		if err != nil {
			return ErrLayers.Wrap(err)
		}
		parts = append(parts, part)
		if canvas.Empty() {
			canvas = full
		}
	}

	number := big.NewInt(int64(token))
	if err := service.metadata.GenerateMetadata(ctx, projectID, number, combination.String(), attributes(projectLayers, indices)); err != nil {
		return ErrLayers.Wrap(err)
	}

	// the images are named by the token id in the random order of the combinations, so a token's traits
	// do not depend on its id; the combination number goes to the dna of the metadata.
	// TODO: add jpeg variant.
	buff := new(bytes.Buffer)
	if err := service.encoder.Encode(buff, imageproc.ComposeOver(canvas.Dx(), canvas.Dy(), parts)); err != nil {
		return ErrLayers.Wrap(err)
	}

	return ErrLayers.Wrap(service.imageStore.SaveImage(ctx, projectID, strconv.Itoa(token)+".png", buff))
}

// layerCache keeps the decoded layer images of one generation, each cropped to its visible pixels: every
// image is read and decoded once instead of once per generated image, and drawing a cropped image touches
// only its visible part. The images are loaded on first use, safely for concurrent workers.
type layerCache struct {
	store   Store
	layers  []Layer
	entries [][]layerCacheEntry
}

type layerCacheEntry struct {
	once  sync.Once
	image image.Image
	full  image.Rectangle
	err   error
}

func newLayerCache(store Store, projectLayers []Layer) *layerCache {
	cache := &layerCache{store: store, layers: projectLayers, entries: make([][]layerCacheEntry, len(projectLayers))}
	for k, layer := range projectLayers {
		cache.entries[k] = make([]layerCacheEntry, len(layer.ImageNames))
	}
	return cache
}

// get returns the cropped image with the number index (from 1) of the k-th layer and the bounds of the whole image.
func (cache *layerCache) get(ctx context.Context, k, index int) (image.Image, image.Rectangle, error) {
	entry := &cache.entries[k][index-1]
	entry.once.Do(func() {
		layer := cache.layers[k]
		reader, err := cache.store.GetLayerImage(ctx, layer.ProjectID, layer.Name, index)
		if err != nil {
			entry.err = err
			return
		}
		decoded, err := png.Decode(reader)
		closeReader(reader)
		if err != nil {
			entry.err = err
			return
		}
		entry.image, entry.full = imageproc.CropToVisible(decoded), decoded.Bounds()
	})
	return entry.image, entry.full, entry.err
}

// closeReader closes a layer image file once it is decoded.
func closeReader(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		_ = closer.Close()
	}
}

// attributes returns the metadata attributes of the image with the given variant numbers: one per layer
// in the layer order, valued by the image file name without the extension, or None for an absent layer.
func attributes(projectLayers []Layer, indices []int) []nft.Attribute {
	attributes := make([]nft.Attribute, len(projectLayers))
	for k, layer := range projectLayers {
		value := NoneName
		if indices[k] <= len(layer.ImageNames) {
			value = strings.TrimSuffix(layer.ImageNames[indices[k]-1], filepath.Ext(layer.ImageNames[indices[k]-1]))
		}
		attributes[k] = nft.Attribute{TraitType: layer.Name, Value: value}
	}
	return attributes
}

// Abort aborts generation.
func (service *Service) Abort(projectID uuid.UUID) error {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	abort, ok := service.abortMap[projectID]
	if !ok {
		return ErrGenerationProcess.Wrap(ErrNoGenerationProcess)
	}
	// the entry goes away at once, so a second abort does not close the channel again.
	close(abort)
	delete(service.abortMap, projectID)

	return nil
}

// CreateLayerImage creates layer image. Layer image is one of the options how the layer will look.
func (service *Service) CreateLayerImage(ctx context.Context, layerID uuid.UUID, fileName string, reader io.Reader) error {
	layer, err := service.layers.Get(ctx, layerID)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	return ErrLayers.Wrap(service.imageStore.Create(ctx, layer.ProjectID, layer.Name, fileName, reader))
}

// DeleteLayerImage deletes layer image and its weight.
func (service *Service) DeleteLayerImage(ctx context.Context, layerID uuid.UUID, fileName string) error {
	layer, err := service.layers.Get(ctx, layerID)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	if err = service.imageStore.Delete(ctx, layer.ProjectID, layer.Name, fileName); err != nil {
		return ErrLayers.Wrap(err)
	}

	return ErrLayers.Wrap(service.layers.DeleteWeight(ctx, layerID, fileName))
}

// UpdateNameOfLayerImage updates name of layer image; the image keeps its weight.
func (service *Service) UpdateNameOfLayerImage(ctx context.Context, layerID uuid.UUID, newFileName, oldFileName string) error {
	layer, err := service.layers.Get(ctx, layerID)
	if err != nil {
		return ErrLayers.Wrap(err)
	}

	if err = service.imageStore.UpdateName(ctx, oldFileName, newFileName, layer.ProjectID, layer.Name); err != nil {
		return ErrLayers.Wrap(err)
	}

	return ErrLayers.Wrap(service.layers.RenameWeight(ctx, layerID, oldFileName, newFileName))
}

// GetImageComponents returns the metadata attributes of the generated image with the given combination number.
func (service *Service) GetImageComponents(ctx context.Context, projectID uuid.UUID, combination *big.Int) ([]nft.Attribute, error) {
	projectLayers, err := service.List(ctx, projectID)
	if err != nil {
		return nil, ErrLayers.Wrap(err)
	}

	return attributes(projectLayers, layerIndices(combination, layerBases(projectLayers))), nil
}
