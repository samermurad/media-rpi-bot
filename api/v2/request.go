package v2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Request interface {
	// SetHeaders()
	Headers() map[string]string
	AddHeader(key, value string)
	Url() string
	// SetUrl(string)

	Method() string
	// SetMethod(string)

	Body() []byte
	// SetBody([]byte)

	SetTimeout(time.Duration)
	Timeout() time.Duration

	Run(ch chan *ResponseChannel)
}
type request struct {
	url     string
	headers map[string]string
	method  string
	body    []byte
	timeout time.Duration
}

func (req *request) SetHeaders(hdrs *http.Header) {
	for k, v := range req.headers {
		hdrs.Set(k, v)
	}
}

func (req *request) Url() string {
	return req.url
}

func (req *request) Method() string {
	return req.method
}

func (req *request) AddHeader(key, value string) {
	req.headers[key] = value
}

func (req *request) Headers() map[string]string {
	return req.headers
}

func (req *request) Body() []byte {
	return req.body
}

func (req *request) SetTimeout(t time.Duration) {
	req.timeout = t
}

func (req *request) Timeout() time.Duration {
	return req.timeout
}

func (req *request) String() string {
	return fmt.Sprintf("Request: %v: %v -- %v", req.method, req.url, req.body)
}

func (req *request) Run(ch chan *ResponseChannel) {
	var endError error = nil
	var endRes *Response = nil
	defer func() {
		ch <- &ResponseChannel{
			Res: endRes,
			Err: endError,
		}
	}()
	_req, err := http.NewRequest(req.Method(), req.Url(), bytes.NewBuffer(req.Body()))
	if err != nil {
		endError = NewError("Failed to Create Request Object")
		return
	}
	req.SetHeaders(&_req.Header)
	fmt.Println("Executing", req)
	client := &http.Client{}
	if req.timeout != 0 {
		client.Timeout = req.timeout
	}

	res, err := client.Do(_req)
	if err != nil {
		endError = err
		return
	} else {
		endRes = &Response{
			Body:   nil,
			Status: res.StatusCode,
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			endError = err
			return
		} else {
			endRes.Body = data
			return
		}
	}
}
