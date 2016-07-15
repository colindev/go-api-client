package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/colindev/go-api-client/test"
)

type Content struct {
	Method        string `json:"method"`
	Path          string `json:"path"`
	ContentType   string `json:"content-type"`
	ContentLength int64  `json:"content-length"`
	FormData      string `json:"form-data"`
}

var c = New("http://127.0.0.1:8000")

func init() {
	c.Replace(test.New().Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		ctn := Content{
			Method:        r.Method,
			Path:          strings.TrimLeft(r.URL.Path, "/"),
			ContentType:   r.Header.Get("Content-Type"),
			ContentLength: r.ContentLength,
		}

		err := r.ParseForm()
		if err != nil {
			ctn.FormData = err.Error()
		} else {
			switch r.Method {
			// https://golang.org/src/net/http/request.go?s=28722:28757#L924
			// NOTE: ACC DELETE Method 需要將 payload 放置 body, 但是 golang http.Request.ParseForm 僅處理 PUT/POST/PATCH 內的 body
			case "DELETE":
				if b, e := ioutil.ReadAll(r.Body); e == nil {
					ctn.FormData = string(b)
				} else {
					ctn.FormData = e.Error()
				}
			default:
				ctn.FormData = r.Form.Encode()
			}
		}

		if b, err := json.Marshal(ctn); err == nil {
			w.Write(b)
		} else {
			w.Write([]byte(err.Error()))
		}
	}))
}

func TestGet(t *testing.T) {

	method := "GET"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := c.Get(path, params)
	if err != nil {
		t.Error(err)
		t.Skip("test content")
	}

	t.Logf("get: %s\n", data)

	var ctn Content
	if err := json.Unmarshal(data, &ctn); err != nil {
		t.Error(err)
	}

	if ctn.Method != method {
		t.Errorf("method expect %s, but %s", method, ctn.Method)
	}

	if ctn.Path != path {
		t.Errorf("path expect %s, but %s", path, ctn.Path)
	}

	if formData := params.Encode(); formData != ctn.FormData {
		t.Errorf("form data encode expect %s, but %s", formData, ctn.FormData)
	}
}

func TestPut(t *testing.T) {

	method := "PUT"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := c.Put(path, params)
	if err != nil {
		t.Error(err)
		t.Skip("test content")
	}

	t.Logf("get: %s\n", data)

	var ctn Content
	if err := json.Unmarshal(data, &ctn); err != nil {
		t.Error(err)
	}

	if ctn.Method != method {
		t.Errorf("method expect %s, but %s", method, ctn.Method)
	}

	if ctn.Path != path {
		t.Errorf("path expect %s, but %s", path, ctn.Path)
	}

	if formData := params.Encode(); formData != ctn.FormData {
		t.Errorf("form data encode expect %s, but %s", formData, ctn.FormData)
	}
}

func TestPost(t *testing.T) {

	method := "POST"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := c.Post(path, params)
	if err != nil {
		t.Error(err)
		t.Skip("test content")
	}

	t.Logf("get: %s\n", data)

	var ctn Content
	if err := json.Unmarshal(data, &ctn); err != nil {
		t.Error(err)
	}

	if ctn.Method != method {
		t.Errorf("method expect %s, but %s", method, ctn.Method)
	}

	if ctn.Path != path {
		t.Errorf("path expect %s, but %s", path, ctn.Path)
	}

	if formData := params.Encode(); formData != ctn.FormData {
		t.Errorf("form data encode expect %s, but %s", formData, ctn.FormData)
	}
}

func TestDelete(t *testing.T) {

	method := "DELETE"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := c.Delete(path, params)
	if err != nil {
		t.Error(err)
		t.Skip("test content")
	}

	t.Logf("get: %s\n", data)

	var ctn Content
	if err := json.Unmarshal(data, &ctn); err != nil {
		t.Error(err)
	}

	if ctn.Method != method {
		t.Errorf("method expect %s, but %s", method, ctn.Method)
	}

	if ctn.Path != path {
		t.Errorf("path expect %s, but %s", path, ctn.Path)
	}

	if formData := params.Encode(); formData != ctn.FormData {
		t.Errorf("form data encode expect %s, but %s", formData, ctn.FormData)
	}
}

func ExampleHttpError() {

	c := New("http://127.0.0.1:8000")
	c.Replace(test.New().Handle(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "test error code", 404)
	}))

	ctn, err := c.Get("error/404", nil)
	fmt.Println(err)
	fmt.Println(string(ctn))
	// Output:
	// http error [404]
	// test error code
}
