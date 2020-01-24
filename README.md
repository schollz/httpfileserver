# httpfileserver
A cache-friendly, gzip-friendly file server to bind the std Golang http


## Example

```golang
func main() {
        http.HandleFunc("/static/", httpfilehandler.New("/static", ".").Handle())
        http.ListenAndServe(":1113", nil)
}
```

## License

MIT
