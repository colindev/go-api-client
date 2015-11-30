# RESTful API Client (by golang)

用 go 實作的 RESTful api client 工具

### Quick Start

```go
client := api.New("http://127.0.0.1/api")

// set header
api.SetHeader("Host", "api.host")

// Get
params, err := url.ParseQuery("a=1&b=2")

// vals = *api.Values
accVals, err := api.Get("hello", params)

fmt.Println(vals.Get("key.key2.key3"))

```

### Trace api

```go
api.Trace(func(req http.Request, b []byte, status int, e error){
    // you can write log here... if e not nil
})
```

