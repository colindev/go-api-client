# RESTful API Client (by golang)

用 go 實作的 RESTful api client 工具

### use custom package

##### Step 1: write in your .ssh/config

```
Host rde-tech.go
    HostName 10.251.39.15
    IdentityFile ~/.ssh/id_rsa.pub
    User dev
    IdentitiesOnly yes
```

##### Step 2: run go get

```sh
go get rde-tech.go/rde-tech/go-api-client.git
```

##### Step 3: write in your go import

```go
import "rde-tech.go/rde-tech/go-acc-client.git"
```

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

