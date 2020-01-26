package httpfileserver

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

type fileServer struct {
	dir        string
	route      string
	middleware middleware
	cache      sync.Map

	optionDisableCache    bool
	optionMaxBytesPerFile int
}

// New returns a new file server that can handle requests for
// files using an in-memory store with gzipping
func New(route, dir string, options ...Option) *fileServer {
	fs := &fileServer{
		dir:                   dir,
		route:                 route,
		optionMaxBytesPerFile: 1000000000,
	}
	for _, o := range options {
		o(fs)
	}
	return fs
}

// Option is the type all options need to adhere to
type Option func(fs *fileServer)

// OptionNoCache disables the caching
func OptionNoCache(disable bool) Option {
	return func(fs *fileServer) {
		fs.optionDisableCache = disable
	}
}

// OptionMaxBytes sets the maximum number of bytes per file to cache
func OptionMaxBytes(optionMaxBytesPerFile int) Option {
	return func(fs *fileServer) {
		fs.optionMaxBytesPerFile = optionMaxBytesPerFile
	}
}

type middleware struct {
	io.Writer
	http.ResponseWriter
	bytesWritten *bytes.Buffer
	numBytes     *int
	overflow     *bool
	maxBytes     int
}

type file struct {
	bytes  []byte
	header http.Header
}

type writeCloser struct {
	*bufio.Writer
}

// Close will close the writer
func (wc *writeCloser) Close() error {
	return wc.Flush()
}

// Write will have the middleware save the bytes
func (m middleware) Write(b []byte) (int, error) {
	if len(b)+*m.numBytes < m.maxBytes {
		n, _ := m.bytesWritten.Write(b)
		*m.numBytes += n
	} else {
		*m.overflow = true
	}
	return m.Writer.Write(b)
}

// Handle gives a handlerfunc for the file server
func (fs *fileServer) Handle() http.HandlerFunc {
	return fs.ServeHTTP
}

// ServeHTTP is the server of the file server
func (fs *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn := http.FileServer(http.Dir(fs.dir)).ServeHTTP
	r.URL.Path = strings.TrimPrefix(r.URL.Path, fs.route)
	doGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
	// check the sync map using the r.URL.Path and return
	// the gzipped or the standard version
	key := r.URL.Path

	if !fs.optionDisableCache {
		if doGzip {
			// load the gzipped cache version if available
			fileint, ok := fs.cache.Load(key + "gzip")
			if ok {
				file := fileint.(file)
				for k := range file.header {
					for _, v := range file.header[k] {
						w.Header().Set(k, v)
					}
				}
				w.Header().Set("Content-Encoding", "gzip")
				w.Write(file.bytes)
				return
			}
		}

		// try to load from cache
		fileint, ok := fs.cache.Load(key)
		if ok {
			file := fileint.(file)
			for k := range file.header {
				if k == "Content-Encoding" {
					continue
				}
				for _, v := range file.header[k] {
					if len(v) == 0 {
						continue
					}
					w.Header().Set(k, v)
				}
			}
			if doGzip {
				w.Header().Set("Content-Encoding", "gzip")
				var wb bytes.Buffer
				wc := gzip.NewWriter(&wb)
				wc.Write(file.bytes)
				wc.Close()
				w.Write(wb.Bytes())
				file.bytes = wb.Bytes()
				fs.cache.Store(key+"gzip", file)
			} else {
				w.Write(file.bytes)
			}
			return
		}
	}

	var wc io.WriteCloser
	if doGzip {
		w.Header().Set("Content-Encoding", "gzip")
		wc = gzip.NewWriter(w)
	} else {
		wc = &writeCloser{bufio.NewWriter(w)}
	}
	defer wc.Close()

	mware := middleware{Writer: wc, ResponseWriter: w, bytesWritten: new(bytes.Buffer), numBytes: new(int), overflow: new(bool), maxBytes: fs.optionMaxBytesPerFile}
	fn(mware, r)

	// extract bytes written and the header and save it as a file
	// to the sync map using the r.URL.Path
	if !fs.optionDisableCache && !*mware.overflow {
		file := file{
			bytes:  mware.bytesWritten.Bytes(),
			header: w.Header(),
		}
		fs.cache.Store(key, file)
	}
}
