package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"rde-tech.go/rde-tech/go-api-client.git"
)

func usage() {
	fmt.Printf("\033[31mUsage: %s %s\033[m\n", os.Args[0], "[GET|POST|PUT|DELETE] [URL] [SUB-URL]")
	os.Exit(2)
}

func title(s string) {
	fmt.Printf("\033[33m%s\n\033[m", s)
}

func content(v interface{}) {
	fmt.Printf("\033[35m%v\n\033[m", v)
}

func main() {

	flag.Parse()

	method := flag.Arg(0)
	base := flag.Arg(1)
	uri := flag.Arg(2)

	if "" == method || "" == base {
		usage()
	}

	client := api.New(base).Trace(func(req http.Request, b []byte, status int, err error) {
		title("Request Method")
		content(req.Method)
		title("Request Proto")
		content(req.Proto)
		title("url.URL String")
		content(req.URL)
		title("tracer status")
		content(status)
		title("tracer body")
		content(string(b))
		title("tracer error")
		fmt.Println(err)
	}).Trace(func(http.Request, []byte, int, error) {
		panic("Test panic")
	})

	methods := map[string]api.ApiHandler{
		"GET":    client.Get,
		"POST":   client.Post,
		"PUT":    client.Put,
		"Delete": client.Delete,
	}

	getter, ok := methods[strings.ToUpper(method)]

	if !ok {
		usage()
	}

	getter(uri, nil)
}
