package metadata

import "generator/pkg/storage"

// MetadataFileType indicates that type of metadata file is json.
const MetadataFileType = "json"

// Config defines configurable variables of metadata.
type Config struct {
	OutputPath  string         `json:"outputPath"`
	FileStorage storage.Config `json:"fileStorage"`
	Bucket      string         `json:"bucket"`
	Metadata    struct {
		Description string `json:"description"`
		ExternalURL string `json:"externalURL"`
		Image       string `json:"image"`
		Name        string `json:"name"`
	} `json:"metadata"`
}
