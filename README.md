# httpfileserver
A cache-friendly, gzip-friendly file server to bind the std Golang http


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

## License

MIT
