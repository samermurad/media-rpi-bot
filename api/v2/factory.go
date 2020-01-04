package v2

import "encoding/json"

type RequestOption func(req *request)

var MethodPost = func(req *request) {
	req.method = "POST"
}
var MethodGet = func(req *request) {
	req.method = "GET"
}

func WithMethod(m string) RequestOption {
	return func(req *request) {
		req.method = m
	}
}

func WithMethodPost() RequestOption {
	return WithMethod("POST")
}

func WithMethodGet() RequestOption {
	return WithMethod("POST")
}

func WithBody(b []byte) RequestOption {
	return func(req *request) {
		req.body = b
	}
}

func WithMarshalableBody(mrshl interface{}) RequestOption {
	return func(req *request) {
		_b, err := json.Marshal(mrshl)
		if err != nil {
			panic(err)
		}
		req.body = _b
	}
}
func NewRequest(path string, opt ...RequestOption) Request {
	req := request{}
	for _, fn := range opt {
		fn(&req)
	}
	return Request(&req)
}
