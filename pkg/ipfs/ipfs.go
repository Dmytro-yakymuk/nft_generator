// Package ipfs adds files to an IPFS node through the Kubo RPC API.
package ipfs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	"github.com/zeebo/errs"
)

// Error is an error class for IPFS errors.
var Error = errs.Class("ipfs error")

// File is a file of a directory added to IPFS. Content is called when the file is sent,
// so a large directory is streamed one file at a time.
type File struct {
	Name    string
	Content func() ([]byte, error)
}

// Client talks to the RPC API of a Kubo node, e.g. http://ipfs:5001.
type Client struct {
	apiURL string
	http   *http.Client
}

// NewClient is a constructor for the Kubo RPC client.
func NewClient(apiURL string) *Client {
	return &Client{apiURL: strings.TrimRight(apiURL, "/"), http: http.DefaultClient}
}

// AddDirectory adds the files as one directory, pins it on the node and returns the CID of the
// directory (CIDv1). The CID depends only on the content, so adding the same files again gives
// the same CID, and a file is addressed as ipfs://<CID>/<name>.
func (client *Client) AddDirectory(ctx context.Context, files []File) (string, error) {
	if client.apiURL == "" {
		return "", Error.New("IPFS node is not configured")
	}

	body, writer := io.Pipe()
	form := multipart.NewWriter(writer)
	go func() {
		writer.CloseWithError(writeFiles(form, files))
	}()

	query := url.Values{
		"wrap-with-directory": {"true"},
		"cid-version":         {"1"},
		"pin":                 {"true"},
		"quieter":             {"true"},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.apiURL+"/api/v0/add?"+query.Encode(), body)
	if err != nil {
		return "", Error.Wrap(err)
	}
	request.Header.Set("Content-Type", form.FormDataContentType())

	response, err := client.http.Do(request)
	if err != nil {
		return "", Error.Wrap(err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return "", Error.New("add: HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(message)))
	}

	// with quieter=true the node writes one JSON line: the wrapping directory.
	var cid string
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		var line struct {
			Hash    string
			Message string
		}
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			return "", Error.Wrap(err)
		}
		if line.Message != "" {
			return "", Error.New("add: %s", line.Message)
		}
		cid = line.Hash
	}
	if err := scanner.Err(); err != nil {
		return "", Error.Wrap(err)
	}
	if cid == "" {
		return "", Error.New("add: the node returned no CID")
	}

	return cid, nil
}

// writeFiles writes the files as parts of the multipart form, then closes it.
func writeFiles(form *multipart.Writer, files []File) error {
	for _, file := range files {
		content, err := file.Content()
		if err != nil {
			return fmt.Errorf("file %s: %w", file.Name, err)
		}

		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, file.Name))
		header.Set("Content-Type", "application/octet-stream")
		part, err := form.CreatePart(header)
		if err != nil {
			return err
		}
		if _, err = part.Write(content); err != nil {
			return err
		}
	}
	return form.Close()
}

// ExportCAR writes the DAG under the CID (e.g. a published directory with all its files) as a CAR
// archive, the format pinning services import with the same root CID.
func (client *Client) ExportCAR(ctx context.Context, cid string, w io.Writer) error {
	if client.apiURL == "" {
		return Error.New("IPFS node is not configured")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.apiURL+"/api/v0/dag/export?"+url.Values{"arg": {cid}}.Encode(), nil)
	if err != nil {
		return Error.Wrap(err)
	}
	response, err := client.http.Do(request)
	if err != nil {
		return Error.Wrap(err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return Error.New("dag export: HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(message)))
	}
	// an error inside the stream arrives as an HTTP trailer after a 200 status.
	if _, err = io.Copy(w, response.Body); err != nil {
		return Error.Wrap(err)
	}
	if message := response.Trailer.Get("X-Stream-Error"); message != "" {
		return Error.New("dag export: %s", message)
	}

	return nil
}
