package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/colindev/go-api-client"
)

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

	var (
		verbose bool
		params  string
		headers = cliHeaders{}
	)
	cli := flag.CommandLine
	cli.Usage = func() {
		fmt.Printf("\033[31mUsage: COMMANDS [GET|POST|PUT|DELETE] [BASE] [URI]\033[m\n")
		cli.PrintDefaults()
		os.Exit(2)

	}
	cli.StringVar(&params, "params", "", "params")
	cli.BoolVar(&verbose, "V", false, "verbose")
	cli.Var(&headers, "H", "headers list")
	cli.Parse(os.Args[1:])

	args := cli.Args()
	if len(args) != 3 {
		cli.Usage()
	}
	method := args[0]
	base := args[1]
	uri := args[2]

	fmt.Println("--- request start")
	defer fmt.Println("--- request end")
	client := api.New(base).Trace(func(req *http.Request, b []byte, status int, err error) {
		fmt.Printf("%s %d %s\n", req.Method, status, req.URL)
		for k, v := range req.Header {
			fmt.Printf("\033[2;33m%s:\033[m %v\n", k, v)
		}

		fmt.Println("\033[2;35m--- body start\033[m")
		fmt.Println(string(b))
		fmt.Println("\033[2;35m--- body end\033[m")
		if err != nil {
			fmt.Printf("\033[2;31m%v\033[m\n", err)
		}
	})

	for k, v := range headers {
		client.SetHeader(k, v)
	}

	methods := map[string]func(string, url.Values) ([]byte, int, error){
		"GET":    client.Get,
		"POST":   client.Post,
		"PUT":    client.Put,
		"Delete": client.Delete,
	}

	getter, ok := methods[strings.ToUpper(method)]

	if !ok {
		cli.Usage()
	}

	query, err := url.ParseQuery(params)
	if err != nil {
		log.Printf("\033[2;31m%v\033[m\n", err)
	}
	getter(uri, query)
}
