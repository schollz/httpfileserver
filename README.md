# httpfileserver

[![travis](https://travis-ci.org/schollz/httpfileserver.svg?branch=master)](https://travis-ci.org/schollz/httpfileserver) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/httpfileserver)](https://goreportcard.com/report/github.com/schollz/httpfileserver) 
[![coverage](https://img.shields.io/badge/coverage-90%25-brightgreen.svg)](https://gocover.io/github.com/schollz/httpfileserver)
[![godocs](https://godoc.org/github.com/schollz/httpfileserver?status.svg)](https://godoc.org/github.com/schollz/httpfileserver) 


This is a drop-in replacement for the stdlib `http.FileServer` that automatically provides gzipping as well as serving from memory instead of from disk. This library wraps the stdlib `http.FileServer` so you get all the benefits of that code, while also providing gzipping and keeping track of bytes and storing served files from memory when they come available.

To use, you can just replace the `http.Fileserver` in your `http.Handle` or `http.HandlerFunc` with `httpfileserver.New(route,directory)`. For example, to serve static assets, you can replace the std lib version

	http.Handle("/static/", http.FileServer(http.Dir(".")))

with this version

	http.Handle("/static/", httpfileserver.New("/static/", "."))

## Usage

In order to serve files from a different directory other than specified by the route, you need to include the route when specifying a file server. For example, if you want to serve `/static` files from a local directory called `/tmp/myassets` then you can specify a new file server in the following:

	http.Handle("/static/", httpfileserver.New("/static/", "/tmp/myassets"))

The route is in the handle as well as the instance of the file server so that it can trim it and then server from the directory as intended.

## Example


```golang
package main

import (
        "net/http"

        "github.com/schollz/httpfileserver"
)

func main() {
	// Any request to /static/somefile.txt will serve the file at the location ./somefile.txt
        http.HandleFunc("/static/", httpfileserver.New("/static/", ".").Handle())
        http.ListenAndServe(":1113", nil)
}
```

## Benchmarks

Using the `example` in this repo I tested both the stdlib and this version for serving a file. This version is about 22% faster (since it is reading from memory) and automatically uses `gzip` when capable.

Using the Go stdlib (`http.Handle("/", http.FileServer(http.Dir(".")))`):

```
$ ab -n 20000 -H "Accept-Encoding: gzip,deflate" http://localhost:1113/main.go # stdlib
...
HTML transferred:       4640000 bytes
Requests per second:    3575.56 [#/sec] (mean)
...
```

Using this library `http.Handle("/new/", httpfileserver.New("/new", "."))`:

```
$ ab -n 20000 -H "Accept-Encoding: gzip,deflate" http://localhost:1113/new/main.go # this lib
...
HTML transferred:       3680000 bytes
Requests per second:    4544.44 [#/sec] (mean)
...
```

## License

MIT
