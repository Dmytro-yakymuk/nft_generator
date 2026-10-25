package images

import "generator/pkg/storage"

// Config defines configurable variables of images.
type Config struct {
	ImagesPath  string         `json:"imagesPath"`
	FileStorage storage.Config `json:"fileStorage"`
	Bucket      string         `json:"bucket"`
}
