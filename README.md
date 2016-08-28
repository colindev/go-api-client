# RESTful API Client (by golang)

用 go 實作的 RESTful api client 工具

[![Build Status](https://travis-ci.org/colindev/go-api-client.svg?branch=master)](https://travis-ci.org/colindev/go-api-client)
[![GoDoc](https://godoc.org/github.com/colindev/go-api-client?status.svg)](https://godoc.org/github.com/colindev/go-api-client)

### Quick Start

```go
client := api.New("http://127.0.0.1/api")

// set header
api.SetHeader("Host", "api.host")

// Get
params, err := url.ParseQuery("a=1&b=2")

// data []byte
data, err := api.Get("hello", params)

```

### Trace api

```go
api.Trace(func(req *http.Request, b []byte, status int, e error){
    // you can write log here... if e not nil
})
```

