package main

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
}

func New(route, dir string) *fileServer {
	return &fileServer{
		dir:   dir,
		route: route,
	}
}

type middleware struct {
	io.Writer
	http.ResponseWriter
	bytesWritten *bytes.Buffer
	numBytes     int
	overflow     bool
}

type file struct {
	bytes  []byte
	header http.Header
}

type writeCloser struct {
	*bufio.Writer
}

func (wc *writeCloser) Close() error {
	return wc.Flush()
}

func (m middleware) Write(b []byte) (int, error) {
	if len(b)+m.numBytes < 1000000 {
		n, _ := m.bytesWritten.Write(b)
		m.numBytes += n
	} else {
		m.overflow = true
	}
	return m.Writer.Write(b)
}

func (fs *fileServer) Handle() http.HandlerFunc {
	return fs.ServeHTTP
}

func (fs *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn := http.FileServer(http.Dir(fs.dir)).ServeHTTP
	r.URL.Path = strings.TrimPrefix(r.URL.Path, fs.route)
	doGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
	// check the sync map using the r.URL.Path and return
	// the gzipped or the standard version
	key := r.URL.Path
	fileint, ok := fs.cache.Load(key)
	if ok {
		file := fileint.(file)
		for k := range file.header {
			for _, v := range file.header[k] {
				if len(v) == 0 {
					continue
				}
				w.Header().Set(k, v)
			}
		}
		if doGzip {
			w.Header().Set("Content-Encoding", "gzip")
			wc := gzip.NewWriter(w)
			defer wc.Close()
			wc.Write(file.bytes)
		} else {
			w.Write(file.bytes)
		}
		return
	}

	var wc io.WriteCloser
	if doGzip {
		wc = gzip.NewWriter(w)
	} else {
		wc = &writeCloser{bufio.NewWriter(w)}
	}
	defer wc.Close()

	mware := middleware{Writer: wc, ResponseWriter: w, bytesWritten: new(bytes.Buffer)}
	fn(mware, r)

	// extract bytes written and the header and save it as a file
	// to the sync map using the r.URL.Path
	if !mware.overflow {
		file := file{
			bytes:  mware.bytesWritten.Bytes(),
			header: w.Header(),
		}
		fs.cache.Store(key, file)
	}
	if doGzip {
		w.Header().Set("Content-Encoding", "gzip")
	}
}

func main() {
	http.HandleFunc("/static/", New("/static", ".").Handle())
	http.ListenAndServe(":1113", nil)
}
