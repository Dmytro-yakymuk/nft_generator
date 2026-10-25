package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeebo/errs"
)

// ensures that Client implements ObjectStorage.
var _ ObjectStorage = (*Client)(nil)

var (
	// MinioError is class for minio errors.
	MinioError = errs.Class("minio")
	// Error is class for remove file storage error.
	Error = errs.Class("remote file storage error")
)

// Config is the setup for a particular client.
type Config struct {
	S3Gateway string `json:"s3Gateway"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Region    string `json:"region" default:"us-east-1"`
}

// Client implements basic S3 Client with minio.
type Client struct {
	API *minio.Client
}

// NewClient creates new Client.
func NewClient(cfg Config) (*Client, error) {
	opts := &minio.Options{
		Creds: credentials.New(
			&credentials.Static{
				Value: credentials.Value{
					AccessKeyID:     cfg.AccessKey,
					SecretAccessKey: cfg.SecretKey,
				},
			},
		),
		Secure: true,
		Region: cfg.Region,
	}

	c, err := minio.New(cfg.S3Gateway, opts)
	if err != nil {
		return &Client{}, MinioError.Wrap(err)
	}

	return &Client{API: c}, nil
}

// Download downloads object from specific bucket and returns it as byte slice.
func (client *Client) Download(ctx context.Context, bucket, objectName string) ([]byte, error) {
	var buffer []byte
	reader, err := client.API.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, MinioError.Wrap(err)
	}
	defer func() { _ = reader.Close() }()

	n, err := reader.Read(buffer[:cap(buffer)])
	if !errors.Is(err, io.EOF) {
		rest, err := ioutil.ReadAll(reader)
		if errors.Is(err, io.EOF) {
			err = nil
		}
		if err != nil {
			return nil, Error.Wrap(err)
		}
		buffer = append(buffer, rest...)
		n = len(buffer)
	}

	buffer = buffer[:n]
	return buffer, nil
}

// Upload uploads provided data into object with specific name into provided bucket.
func (client *Client) Upload(ctx context.Context, bucket, objectName string, data []byte) error {
	_, err := client.API.PutObject(ctx, bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return MinioError.Wrap(err)
	}

	return nil
}

// Delete deletes object by object key in specific bucket.
func (client *Client) Delete(ctx context.Context, bucket, objectName string) error {
	return MinioError.Wrap(client.API.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{}))
}

// ListKeys return list of object keys from requested bucket.
func (client *Client) ListKeys(ctx context.Context, bucket string, options minio.ListObjectsOptions) (keys []string) {
	for object := range client.API.ListObjects(ctx, bucket, options) {
		keys = append(keys, object.Key)
	}

	return
}
