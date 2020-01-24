package main

import (
	"net/http"

	"github.com/schollz/httpfileserver"
)

func main() {
	http.HandleFunc("/static/", httpfileserver.New("/static", ".").Handle())
	http.ListenAndServe(":1113", nil)
}
