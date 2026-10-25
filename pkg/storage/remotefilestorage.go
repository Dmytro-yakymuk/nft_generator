package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
)

// ObjectStorage interface to call s3 with minio.
type ObjectStorage interface {
	Upload(ctx context.Context, bucket, objectName string, data []byte) error
	Download(ctx context.Context, bucket, objectName string) ([]byte, error)
	Delete(ctx context.Context, bucket, objectName string) error
	ListKeys(ctx context.Context, bucket string, options minio.ListObjectsOptions) []string
}
