# httpfileserver
A cache-friendly, gzip-friendly file server to bind the std Golang http

```
BenchmarkServer-8       2020/01/24 17:50:34 http: Accept error: accept tcp 127.0.0.1:35747: accept4: too many open files; retrying in 5ms
2020/01/24 17:50:34 Get http://127.0.0.1:35747/README.md: dial tcp 127.0.0.1:35747: socket: too many open files
exit status 1
```

