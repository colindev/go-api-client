package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ApiHandler func(string, url.Values) (*Values, error)
type ApiTracer func(http.Request, []byte, error)

type headers map[string]string

type api struct {
	client   *http.Client
	base_url string
	headers  headers
	tracers  []ApiTracer
}

func New(base string) *api {

	base = strings.TrimRight(base, "/")

	return &api{
		client:   &http.Client{},
		base_url: base,
		headers:  make(headers),
	}
}

// 設定共用檔頭
func (a *api) SetHeader(name, value string) *api {
	a.headers[name] = value

	return a
}

// 注入追蹤程序
func (a *api) Trace(tc ApiTracer) *api {
	a.tracers = append(a.tracers, tc)

	return a
}

// GET
func (a *api) Get(uri string, params url.Values) (v *Values, err error) {

	uri = resolveUri(uri)
	req, err := http.NewRequest("GET", a.base_url+"/"+uri+"?"+params.Encode(), nil)

	return resolveRequest(a, req, err)
}

// POST
func (a *api) Post(uri string, params url.Values) (v *Values, err error) {

	uri = resolveUri(uri)
	req, err := http.NewRequest("POST", a.base_url+"/"+uri, bytes.NewBufferString(params.Encode()))

	return resolveRequest(a, req, err)
}

// PUT
func (a *api) Put(uri string, params url.Values) (v *Values, err error) {

	uri = resolveUri(uri)
	req, err := http.NewRequest("PUT", a.base_url+"/"+uri, bytes.NewBufferString(params.Encode()))

	return resolveRequest(a, req, err)
}

// DELETE
func (a *api) Delete(uri string, params url.Values) (v *Values, err error) {

	uri = resolveUri(uri)
	req, err := http.NewRequest("DELETE", a.base_url+"/"+uri, bytes.NewBufferString(params.Encode()))

	return resolveRequest(a, req, err)
}

func resolveUri(s string) string {
	return strings.TrimLeft(s, "/")
}

func resolveHeaders(req *http.Request, headers headers) {
	for name, value := range headers {
		req.Header.Set(name, value)
	}
}

func resolveRequest(a *api, req *http.Request, e error) (v *Values, err error) {

	if e != nil {
		err = fmt.Errorf("new request error( %s )", err)
		return
	}

	resolveHeaders(req, a.headers)

	res, err := a.client.Do(req)
	if err != nil {
		err = fmt.Errorf("request send error( %s )", err)
		return
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	for _, tc := range a.tracers {
		tc(*req, content, err)
	}
	if err != nil {
		err = fmt.Errorf("response body read error( %s )", err)
		return
	}

	// response body to json
	data := make(map[string]interface{})
	if e := json.Unmarshal(content, &data); e != nil {
		err = fmt.Errorf("json.Unmarshal error( %s )", e)
		return
	}

	// json 資料內必須要有 result 這個 key, 且其值必須為 "ok"
	// 大小寫不嚴格檢查
	result, ok := data["result"]
	if !ok {
		err = errors.New("response body miss the key [result]")
		return
	}

	x, ok := result.(string)
	if !ok {
		err = fmt.Errorf("( %v ) is not string", result)
		return
	}
	if "ok" != strings.ToLower(x) {
		err = fmt.Errorf("result not ok ( %s )", result)
		return
	}

	ret, ok := data["ret"]
	if !ok {
		err = errors.New("回應資料有誤")
	}

	v = &Values{ret}
	return
}
