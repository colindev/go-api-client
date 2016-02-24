package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/colin1124x/go-api-client"
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

type cliHeaders map[string]string

func (ch *cliHeaders) String() string {
	return fmt.Sprintf("%s", *ch)
}

func (ch *cliHeaders) Set(v string) (e error) {
	arr := strings.SplitN(v, ":", 2)
	if 2 != len(arr) {
		e = fmt.Errorf("header mush be \"key:value\"")
		return
	}

	(*ch)[strings.Trim(arr[0], " ")] = strings.Trim(arr[1], " ")
	return
}

func main() {

	var headers cliHeaders = map[string]string{}
	flag.Var(&headers, "H", "headers list")
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
		title("Request Header")
		content(req.Header)
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
	}).Trace(func(http.Request, []byte, int, error) {
		title("after panic")
	})

	for k, v := range headers {
		client.SetHeader(k, v)
	}

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
