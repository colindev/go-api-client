package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ApiHandler func(string, url.Values) ([]byte, error)
type Tracer func(*http.Request, []byte, int, error)
type Resolver interface {
	Do(*http.Request) (*http.Response, error)
}

type headers map[string]string

type Client interface {
	Replace(Resolver) Client
	SetHeader(name, value string) Client
	Trace(Tracer) Client
	Get(string, url.Values) ([]byte, error)
	Post(string, url.Values) ([]byte, error)
	Put(string, url.Values) ([]byte, error)
	Delete(string, url.Values) ([]byte, error)
}

type client struct {
	Resolver
	base_url string
	headers  headers
	tracers  []Tracer
}

func New(base string) Client {

	base = strings.TrimRight(base, "/")

	return &client{
		Resolver: &http.Client{},
		base_url: base,
		headers:  make(headers),
	}
}

func (c *client) Replace(r Resolver) Client {
	c.Resolver = r
	return c
}

// 設定共用檔頭
func (c *client) SetHeader(name, value string) Client {

	name = strings.Title(name)

	c.headers[name] = value

	return c
}

// 注入追蹤程序
func (c *client) Trace(tc Tracer) Client {
	c.tracers = append(c.tracers, func(r *http.Request, b []byte, i int, e error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		tc(r, b, i, e)
	})

	return c
}

// GET
func (c *client) Get(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(c.base_url, uri)
	if params != nil {
		resource += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", resource, nil)

	return resolveRequest(c, req, err)
}

// POST
func (c *client) Post(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(c.base_url, uri)

	req, err := http.NewRequest("POST", resource, strings.NewReader(params.Encode()))

	return resolveRequest(c, req, err)
}

// PUT
func (c *client) Put(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(c.base_url, uri)

	req, err := http.NewRequest("PUT", resource, strings.NewReader(params.Encode()))

	return resolveRequest(c, req, err)
}

// DELETE
func (c *client) Delete(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(c.base_url, uri)

	req, err := http.NewRequest("DELETE", resource, strings.NewReader(params.Encode()))

	return resolveRequest(c, req, err)
}

func resolveUrl(base, s string) string {

	base = strings.SplitN(base, "?", 2)[0]

	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(s, "/")
}

func resolveHeaders(req *http.Request, headers headers) {
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	if host, ok := headers["Host"]; ok {
		req.Host = host
	}
}

func resolveTracers(tcs []Tracer, req *http.Request, ctn []byte, sc int, err error) {
	for _, tc := range tcs {
		tc(req, ctn, sc, err)
	}
}

func resolveRequest(c *client, req *http.Request, e error) (ctn []byte, err error) {
	var (
		tracers []Tracer      = c.tracers
		request *http.Request = req
		status  int
	)
	defer func() {
		resolveTracers(tracers, request, ctn, status, err)
	}()

	if e != nil {
		err = fmt.Errorf("new request error( %s )", err)
		return
	}

	resolveHeaders(req, c.headers)

	switch req.Method {
	case "PUT", "POST", "DELETE":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	res, err := c.Resolver.Do(req)
	if err != nil {
		err = fmt.Errorf("request send error( %s )", err)
		return
	}
	defer res.Body.Close()

	ctn, err = ioutil.ReadAll(res.Body)

	if err != nil {
		err = fmt.Errorf("response body read error( %s )", err)
		return
	}

	status = res.StatusCode
	if status != http.StatusOK {
		err = fmt.Errorf("http error [%d]", status)
	}

	return
}
