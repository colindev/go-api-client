package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

type echo struct {
	ln net.Listener
}

func (s *echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var res string
	switch r.URL.Path {
	case "/sum":
		a, err := strconv.Atoi(q.Get("a"))
		if err != nil {
			res = err.Error()
			break
		}
		b, err := strconv.Atoi(q.Get("b"))
		if err != nil {
			res = err.Error()
			break
		}

		res = strconv.Itoa(a + b)
	default:
		res = r.URL.Path
	}

	w.Write([]byte(res))
}

func startEcheServer(t *testing.T, addr string) (ln net.Listener, closed chan bool, err error) {
	ln, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}

	closed = make(chan bool)
	go func() {
		http.Serve(ln, &echo{ln})
		close(closed)
	}()
	return ln, closed, err
}

var (
	port string = "8000"
)

func TestResponse(t *testing.T) {
	ln, closed, err := startEcheServer(t, ":"+port)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	fmt.Println("server start")

	client := New("http://127.0.0.1:" + port)

	var (
		data []byte
		e    error
	)

	data, e = client.Get("ok", nil)
	if e != nil {
		t.Error(e)
	} else if string(data) != "/ok" {
		t.Errorf("expect: /ok but %s", data)
	}
	params, e := url.ParseQuery("a=1&b=2")
	if e != nil {
		t.Error(e)
		return
	}
	data, e = client.Get("sum", params)
	if e != nil {
		t.Error(e)
	} else if string(data) != "3" {
		t.Errorf("expect: 3 but %s", data)
	}

	ln.Close()
	<-closed
	fmt.Println("server closed")
}

func TestTracer(t *testing.T) {
	client := New("http://127.0.0.1")
	client.Trace(func(req *http.Request, ctn []byte, sc int, err error) {
		if req.URL.Path != "/x" {
			t.Error("url must be '/x', but:", req.URL.Path)
		}

		if string(ctn) != "x" {
			t.Errorf("content must be 'x', but: %s", ctn)
		}

		if sc != 200 {
			t.Error("status code must be 200, but:", sc)
		}

		if err == nil || err.Error() != "x" {
			t.Errorf("error must be error('x'), but %v", err)
		}
	})

	req, err := http.NewRequest("get", "http://127.0.0.1/x", nil)
	if err != nil {
		t.Error(err)
		return
	}
	client.FireTracers(req, []byte("x"), 200, errors.New("x"))
}
