// Package filebase pins IPFS content in Filebase: a CAR archive uploaded to an IPFS bucket through the
// S3 API with the metadata import=car is unpacked and pinned under the CID of its root, so the files stay
// available in IPFS without the platform's node.
package filebase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/minio/minio-go/v7/pkg/signer"
	"github.com/zeebo/errs"
)

// Error is an error class for Filebase errors.
var Error = errs.Class("filebase error")

// Service is the name of the pinning service shown to the owner.
const Service = "Filebase"

// Config holds the IPFS bucket and the S3 access key of a Filebase account.
type Config struct {
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// Enabled reports whether the account is configured.
func (config Config) Enabled() bool {
	return config.Bucket != "" && config.AccessKey != "" && config.SecretKey != ""
}

// Client uploads CAR archives to a Filebase IPFS bucket.
type Client struct {
	config Config
	http   *http.Client
}

// NewClient is a constructor for the Filebase client.
func NewClient(config Config) *Client {
	if config.Endpoint == "" {
		config.Endpoint = "https://s3.filebase.com"
	}
	config.Endpoint = strings.TrimRight(config.Endpoint, "/")

	return &Client{config: config, http: http.DefaultClient}
}

// PinCAR uploads the CAR archive whose root is cid as the object <cid>.car and checks that Filebase
// pinned the same CID (it returns the root in the x-amz-meta-cid header). The archive is read twice:
// first for the payload hash the request is signed with, then as the body.
func (client *Client) PinCAR(ctx context.Context, cid string, car *os.File) error {
	hash := sha256.New()
	size, err := io.Copy(hash, car)
	if err != nil {
		return Error.Wrap(err)
	}
	if _, err = car.Seek(0, io.SeekStart); err != nil {
		return Error.Wrap(err)
	}

	object := client.config.Endpoint + "/" + url.PathEscape(client.config.Bucket) + "/" + url.PathEscape(cid) + ".car"
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, object, io.NopCloser(car))
	if err != nil {
		return Error.Wrap(err)
	}
	request.ContentLength = size
	request.Header.Set("Content-Type", "application/vnd.ipld.car")
	request.Header.Set("X-Amz-Meta-Import", "car")
	request.Header.Set("X-Amz-Content-Sha256", hex.EncodeToString(hash.Sum(nil)))
	request = signer.SignV4(*request, client.config.AccessKey, client.config.SecretKey, "", "us-east-1")

	response, err := client.http.Do(request)
	if err != nil {
		return Error.Wrap(err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return Error.New("upload %s.car: HTTP %d: %s", cid, response.StatusCode, strings.TrimSpace(string(message)))
	}
	if pinned := response.Header.Get("X-Amz-Meta-Cid"); pinned != cid {
		return Error.New("pinned %q instead of %s", pinned, cid)
	}

	return nil
}
