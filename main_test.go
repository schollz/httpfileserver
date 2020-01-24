package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(New("/", ".").Handle())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/README.md")
	if err != nil {
		t.Errorf("%s", err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s", greeting)

}
