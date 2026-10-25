package httpserver

import (
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// precompressed serves the files of a directory like http.FileServer, but sends the compressed copy made by
// the build (name.br or name.gz) when the browser accepts its encoding. The build names its files by their
// content hash, so the browser may keep the files under /dist/ for a year.
type precompressed struct {
	dir   string
	files http.Handler
}

// encodings are the compressed copies in the order of preference.
var encodings = []struct{ name, extension string }{{"br", ".br"}, {"gzip", ".gz"}}

func newPrecompressed(dir string) http.Handler {
	return &precompressed{dir: dir, files: http.FileServer(http.Dir(dir))}
}

func (handler *precompressed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := path.Clean("/" + r.URL.Path)
	if strings.HasPrefix(name, "/dist/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	}
	w.Header().Add("Vary", "Accept-Encoding")

	for _, encoding := range encodings {
		if !acceptsEncoding(r.Header.Get("Accept-Encoding"), encoding.name) {
			continue
		}
		file := filepath.Join(handler.dir, filepath.FromSlash(name)+encoding.extension)
		if info, err := os.Stat(file); err != nil || info.IsDir() {
			continue
		}
		// the type is that of the original file, not of .br or .gz.
		if contentType := mime.TypeByExtension(path.Ext(name)); contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.Header().Set("Content-Encoding", encoding.name)
		http.ServeFile(w, r, file)
		return
	}

	handler.files.ServeHTTP(w, r)
}

// acceptsEncoding reports whether the Accept-Encoding header allows the encoding: it is listed (or * is)
// with a quality above zero.
func acceptsEncoding(header, encoding string) bool {
	for _, part := range strings.Split(header, ",") {
		fields := strings.Split(part, ";")
		token := strings.TrimSpace(fields[0])
		if token != encoding && token != "*" {
			continue
		}
		quality := 1.0
		for _, parameter := range fields[1:] {
			if value := strings.TrimSpace(parameter); strings.HasPrefix(value, "q=") {
				if q, err := strconv.ParseFloat(strings.TrimPrefix(value, "q="), 64); err == nil {
					quality = q
				}
			}
		}
		return quality > 0
	}
	return false
}
