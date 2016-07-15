package test

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHandle(t *testing.T) {
	c := New().Handle(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Write([]byte("a"))
		case "POST":
			w.Write([]byte("b"))
		case "PUT":
			w.Write([]byte("c"))
		case "DELETE":
			w.Write([]byte("d"))
		default:
			w.Write([]byte("x"))
		}
	})

	var (
		err error
		req *http.Request
		res *http.Response
	)

	for m, b := range map[string]string{"GET": "a", "POST": "b", "PUT": "c", "DELETE": "d"} {

		req, err = http.NewRequest(m, "http://127.0.0.1", nil)
		if err != nil {
			t.Error(err)
			continue
		}
		res, err = c.Do(req)
		if err != nil {
			t.Error(err)
			continue
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			continue
		}

		if string(body) != b {
			t.Errorf("body: expect %s, but %s", b, body)
		}

		if res.StatusCode != 500 {
			t.Errorf("status code: expect %d, but %d", 500, res.StatusCode)
		}
	}
}
