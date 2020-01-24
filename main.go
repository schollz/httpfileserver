package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
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
	m.bytesWritten.Write(b)
	return m.Writer.Write(b)
}

func (fs *fileServer) Handle() http.HandlerFunc {
	fn := http.FileServer(http.Dir(fs.dir)).ServeHTTP
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, fs.route)
		doGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		// TODO
		// check the sync map using the r.URL.Path and return
		// the gzipped or the standard version

		var wc io.WriteCloser
		if doGzip {
			w.Header().Set("Content-Encoding", "gzip")
			wc = gzip.NewWriter(w)
		} else {
			wc = &writeCloser{bufio.NewWriter(w)}
		}
		defer wc.Close()

		gzr := middleware{Writer: wc, ResponseWriter: w, bytesWritten: new(bytes.Buffer)}
		fn(gzr, r)

		// TODO
		// extract bytes written and the header and save it as a file
		// to the sync map using the r.URL.Path
	}
}

func main() {
	log.Println("running on 1113")
	http.HandleFunc("/static/", New("/static", ".").Handle())
	http.ListenAndServe(":1113", nil)
}
