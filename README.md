# httpfileserver

This is a drop-in replacement for the Golang stdlib `http.FileServer`. Instead of writing

	http.Handle("/", http.FileServer(http.Dir(".")))

you can just write

	http.Handle("/", httpfileserver.New("/", "."))


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
