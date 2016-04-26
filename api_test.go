package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type Content struct {
	Method        string `json:"method"`
	Path          string `json:"path"`
	ContentType   string `json:"content-type"`
	ContentLength int64  `json:"content-length"`
	FormData      string `json:"form-data"`
}

func init() {

	router := mux.NewRouter()
	router.HandleFunc("/{path:.+}", func(w http.ResponseWriter, r *http.Request) {
		var content Content
		content.Method = r.Method
		content.Path = strings.TrimLeft(r.URL.Path, "/")
		content.ContentLength = r.ContentLength
		content.ContentType = r.Header.Get("Content-Type")

		err := r.ParseForm()
		if err != nil {
			content.FormData = err.Error()
		} else {
			content.FormData = r.Form.Encode()
		}

		if b, err := json.Marshal(content); err == nil {
			w.Write(b)
		} else {
			w.Write([]byte(err.Error()))
		}
	})

	go http.ListenAndServe(":8000", router)
	time.Sleep(time.Nanosecond * 5000)
}

func TestGet(t *testing.T) {

	client := New("http://127.0.0.1:8000")

	method := "GET"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := client.Get(path, params)
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

	client := New("http://127.0.0.1:8000")

	method := "PUT"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := client.Put(path, params)
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

	client := New("http://127.0.0.1:8000")

	method := "POST"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := client.Post(path, params)
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

	client := New("http://127.0.0.1:8000")

	method := "DELETE"
	path := "ok"
	params := url.Values{}
	params.Set("a", "1")
	params.Set("b", "2")
	params.Set("c[0][x]", "3")
	data, err := client.Delete(path, params)
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

	// https://golang.org/src/net/http/request.go?s=28722:28757#L924
	// BUG ACC DELETE Method 需要將 payload 放置 body, 但是 golang http.Request.ParseForm 僅處理 PUT/POST/PATCH 內的 body
	//if formData := params.Encode(); formData != ctn.FormData {
	//	t.Errorf("form data encode expect %s, but %s", formData, ctn.FormData)
	//}
}
