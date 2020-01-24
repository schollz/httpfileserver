# httpfileserver

This is a drop-in replacement for the Golang stdlib `http.FileServer` that serves from memory instead of from disk, as well as provides automatic gzipping when requested.

To use, you can just replace your `http.Handle` or `http.HandlerFunc`, e.g. you can take the stdlib version:

	http.Handle("/", http.FileServer(http.Dir(".")))

and replace it with this library:

	http.Handle("/", httpfileserver.New("/", "."))

This library essentially wraps the stdlib `http.FileServer` so you get all the benefits of that library, while also providing gzipping and keeping track of bytes and storing served files from memory when they come available.

## Example


```golang
package main

import (
        "net/http"

        "github.com/schollz/httpfileserver"
)

func main() {
	// Any request to /static/somefile.txt will serve the file at the location ./somefile.txt
        http.HandleFunc("/static/", httpfileserver.New("/static", ".").Handle())
        http.ListenAndServe(":1113", nil)
}
```

## Benchmarks

Using the `example` I tested both the stdlib and this version for serving a file. This version is about 22% faster (since it is reading from memory) and automatically uses `gzip` when capable.

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
