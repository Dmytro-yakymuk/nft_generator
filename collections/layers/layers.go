package layers

import (
	"context"
	"image/png"
	"io"
	"sync"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

var (
	// ErrNoLayer indicated that layer does not exist.
	ErrNoLayer = errs.Class("layer does not exist")
	// ErrNoLayerFile indicated that layer file does not exist.
	ErrNoLayerFile = errs.Class("layer file does not exist")
	// ErrGenerationProcess indicates that there was an error related to generation process.
	ErrGenerationProcess = errs.Class("generation process error")
	// ErrNoGenerationProcess indicates that generation progress does not exists.
	ErrNoGenerationProcess = errs.New("generation process does not exists")
	// ErrInvalidDimension indicates that image dimension is invalid.
	ErrInvalidDimension = errs.New("invalid image dimension")
	// ErrInvalidWeight indicates that a variant weight is out of range or names no image of the layer.
	ErrInvalidWeight = errs.Class("invalid weight")
	// ErrRareCombinations indicates that the weights leave too few likely combinations for the collection size.
	ErrRareCombinations = errs.Class("not enough distinct combinations for these weights")
)

// DB exposes access to layers db.
//
// architecture: DB
type DB interface {
	// Create creates a layers and writes to the database.
	Create(ctx context.Context, layer Layer) error
	// List returns all project layers from the database.
	List(ctx context.Context, projectID uuid.UUID) ([]Layer, error)
	// Get returns layer by id from the database.
	Get(ctx context.Context, id uuid.UUID) (Layer, error)
	// UpdatePriorities updates priorities of layers in the database.
	UpdatePriorities(ctx context.Context, layer []Layer) error
	// Delete deletes a layers from the database.
	Delete(ctx context.Context, id uuid.UUID) error
	// ListWeights returns the set weights of the images of every project layer, by layer id and file name.
	ListWeights(ctx context.Context, projectID uuid.UUID) (map[uuid.UUID]map[string]int, error)
	// UpdateWeights replaces the weight of the layer being absent and the weights of its images.
	UpdateWeights(ctx context.Context, layerID uuid.UUID, noneWeight int, weights map[string]int) error
	// RenameWeight moves the weight of a renamed image to its new name.
	RenameWeight(ctx context.Context, layerID uuid.UUID, oldName, newName string) error
	// DeleteWeight deletes the weight of a deleted image.
	DeleteWeight(ctx context.Context, layerID uuid.UUID, name string) error
}

// Store exposes access to layer images storage.
type Store interface {
	// MakeDirectory creates directory in the file system by specific path.
	MakeDirectory(ctx context.Context, path string) error
	// SaveImage saves generated image.
	// TODO: Remove one, just add a different path.
	SaveImage(ctx context.Context, projectID uuid.UUID, imageName string, reader io.Reader) error
	// Create creates layer image.
	Create(ctx context.Context, projectID uuid.UUID, layerName, layerImgName string, reader io.Reader) error
	// GetImageName returns image name by sequence number of image.
	GetImageName(ctx context.Context, projectID uuid.UUID, layerName string, imgSequenceNumber int) (string, error)
	// GetLayerImage returns layer image by sequence number of image.
	GetLayerImage(ctx context.Context, projectID uuid.UUID, layerName string, imgSequenceNumber int) (io.Reader, error)
	// ListLayerImages returns all names of layer images.
	ListLayerImages(ctx context.Context, projectID uuid.UUID, layerName string) ([]string, error)
	// UpdateName updates name of layer image.
	UpdateName(ctx context.Context, oldName, newName string, projectID uuid.UUID, layerName string) error
	// Delete deletes layer image.
	Delete(ctx context.Context, projectID uuid.UUID, layerName string, layerImageName string) error
	// DeleteFolder deletes layer folder with content inside.
	DeleteFolder(ctx context.Context, projectID uuid.UUID, layerName string) error
}

// Config defines configurable values for layers.
type Config struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	// Workers is the number of images generated at once; zero means one per processor core.
	Workers int `json:"workers"`
	// PNGCompression is the zlib level of the generated images: "default", "best-speed",
	// "best-compression" or "none". Encoding takes most of the generation time, and best-speed
	// encodes several times faster for about a third larger files.
	PNGCompression string `json:"pngCompression"`
}

// pngLevels maps the PNGCompression setting to the compression levels of image/png.
var pngLevels = map[string]png.CompressionLevel{
	"":                 png.DefaultCompression,
	"default":          png.DefaultCompression,
	"best-speed":       png.BestSpeed,
	"best-compression": png.BestCompression,
	"none":             png.NoCompression,
}

// PNGEncoder returns the encoder of the generated images for the PNGCompression setting.
func (config Config) PNGEncoder() (*png.Encoder, error) {
	level, ok := pngLevels[config.PNGCompression]
	if !ok {
		return nil, ErrLayers.New("unknown PNG compression %q: use default, best-speed, best-compression or none", config.PNGCompression)
	}

	return &png.Encoder{CompressionLevel: level, BufferPool: &pngBuffers{}}, nil
}

// pngBuffers lets the workers reuse the encoder buffers of finished images.
type pngBuffers struct {
	pool sync.Pool
}

// Get returns a free encoder buffer or nil, which makes the encoder allocate one.
func (buffers *pngBuffers) Get() *png.EncoderBuffer {
	buffer, _ := buffers.pool.Get().(*png.EncoderBuffer)
	return buffer
}

// Put returns the buffer of an encoded image to the pool.
func (buffers *pngBuffers) Put(buffer *png.EncoderBuffer) {
	buffers.pool.Put(buffer)
}

// Layer defines layer entity.
type Layer struct {
	ID               uuid.UUID `json:"id"`
	ProjectID        uuid.UUID `json:"projectId"`
	Name             string    `json:"name"`
	Priority         int       `json:"priority"`
	LayerImagesCount int       `json:"layerImagesCount"`
	// NoneWeight is the weight of the layer being absent from an image; zero makes the layer mandatory.
	NoneWeight int `json:"noneWeight"`
	// Weights holds the weights of the layer images by file name; an image without one has DefaultWeight.
	Weights map[string]int `json:"weights"`
	// ImageNames lists the layer images in the order of their numbers.
	ImageNames []string `json:"-"`
}

// NoneName is the attribute value of a layer that is absent from an image.
const NoneName = "None"

// Weight returns the weight of the layer image with the given file name.
func (layer Layer) Weight(name string) int {
	if weight, ok := layer.Weights[name]; ok {
		return weight
	}
	return DefaultWeight
}
