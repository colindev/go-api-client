package test

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
	Handle(func(http.ResponseWriter, *http.Request)) Client
}

type client struct {
	handler func(http.ResponseWriter, *http.Request)
}

func New() Client {
	return &client{}
}

func (c *client) Handle(handler func(w http.ResponseWriter, r *http.Request)) Client {

	c.handler = handler
	return c
}

func (c *client) Do(req *http.Request) (*http.Response, error) {

	res := newResponse()
	c.handler(res, req)
	if res.code <= 0 {
		res.code = 500
	}

	return http.ReadResponse(res.GetReader(), req)
}

type response struct {
	code   int
	header http.Header
	w      *bytes.Buffer
}

func newResponse() *response {
	return &response{
		w:      bytes.NewBuffer([]byte{}),
		header: make(map[string][]string),
	}
}

func (res *response) Header() http.Header {
	return res.header
}

func (res *response) Write(b []byte) (int, error) {
	return res.w.Write(b)
}

func (res *response) WriteHeader(i int) {
	res.code = i
}

func (res *response) GetReader() *bufio.Reader {

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("HTTP/1.1 %d X\n", res.code))
	res.header.Write(buf)
	buf.WriteString(fmt.Sprintf("Content-Length: %d\n", res.w.Len()))
	buf.WriteRune('\n')
	res.w.WriteTo(buf)

	return bufio.NewReader(buf)
}
