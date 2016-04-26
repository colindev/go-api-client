package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ApiHandler func(string, url.Values) ([]byte, error)
type ApiTracer func(*http.Request, []byte, int, error)

type headers map[string]string

type Api struct {
	client   *http.Client
	base_url string
	headers  headers
	tracers  []ApiTracer
}

func New(base string) *Api {

	base = strings.TrimRight(base, "/")

	return &Api{
		client:   &http.Client{},
		base_url: base,
		headers:  make(headers),
	}
}

// 設定共用檔頭
func (a *Api) SetHeader(name, value string) *Api {

	name = strings.Title(name)

	a.headers[name] = value

	return a
}

// 注入追蹤程序
func (a *Api) Trace(tc ApiTracer) *Api {
	a.tracers = append(a.tracers, func(r *http.Request, b []byte, i int, e error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		tc(r, b, i, e)
	})

	return a
}

// GET
func (a *Api) Get(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(a.base_url, uri)
	if params != nil {
		resource += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", resource, nil)

	return resolveRequest(a, req, err)
}

// POST
func (a *Api) Post(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(a.base_url, uri)

	req, err := http.NewRequest("POST", resource, strings.NewReader(params.Encode()))

	return resolveRequest(a, req, err)
}

// PUT
func (a *Api) Put(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(a.base_url, uri)

	req, err := http.NewRequest("PUT", resource, strings.NewReader(params.Encode()))

	return resolveRequest(a, req, err)
}

// DELETE
func (a *Api) Delete(uri string, params url.Values) ([]byte, error) {

	resource := resolveUrl(a.base_url, uri)

	req, err := http.NewRequest("DELETE", resource, strings.NewReader(params.Encode()))

	return resolveRequest(a, req, err)
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

func resolveTracers(tcs []ApiTracer, req *http.Request, ctn []byte, sc int, err error) {
	for _, tc := range tcs {
		tc(req, ctn, sc, err)
	}
}

func resolveRequest(a *Api, req *http.Request, e error) (ctn []byte, err error) {
	var (
		tracers []ApiTracer   = a.tracers
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

	resolveHeaders(req, a.headers)

	switch req.Method {
	case "PUT", "POST", "DELETE":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	res, err := a.client.Do(req)
	if err != nil {
		err = fmt.Errorf("request send error( %s )", err)
		return
	}
	defer res.Body.Close()
	status = res.StatusCode

	ctn, err = ioutil.ReadAll(res.Body)

	if err != nil {
		err = fmt.Errorf("response body read error( %s )", err)
		return
	}

	return
}
