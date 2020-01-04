package v2

import (
	"encoding/json"
)

type RequestBuilder interface {
	Post() RequestBuilder
	Get() RequestBuilder
	SetHeaders(map[string]string) RequestBuilder
	AddHeader(key, value string) RequestBuilder
	SetBody([]byte) RequestBuilder
	MarshalBody(interface{}) RequestBuilder
	Build() Request
	AppendUrl(string) RequestBuilder
}

type reqBuild struct {
	req request
}

func NewBuilder(url string) RequestBuilder {
	reb := reqBuild{
		req: request{
			url: url,
		},
	}
	return reb
}

func (reqB reqBuild) Post() RequestBuilder {
	reqB.req.method = "POST"
	return reqB
}

func (reqB reqBuild) Get() RequestBuilder {
	reqB.req.method = "GET"
	return reqB
}
func (reqB reqBuild) SetHeaders(dict map[string]string) RequestBuilder {
	// reqB.req.SetHeaders
	if reqB.req.headers == nil {
		reqB.req.headers = make(map[string]string)
	}
	for k, v := range dict {
		reqB.req.headers[k] = v
	}

	return reqB
}

func (reqB reqBuild) AddHeader(key, value string) RequestBuilder {
	if reqB.req.headers == nil {
		reqB.req.headers = make(map[string]string)
	}
	reqB.req.headers[key] = value
	return reqB
}

func (reqB reqBuild) SetBody(data []byte) RequestBuilder {
	reqB.req.body = data
	return reqB
}

func (reqB reqBuild) MarshalBody(data interface{}) RequestBuilder {
	_b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	reqB.req.body = _b
	return reqB
}

func (reqB reqBuild) AppendUrl(url string) RequestBuilder {
	reqB.req.url += url
	return reqB
}
func (reqB reqBuild) Build() Request {
	newR := reqB.req
	return &newR
}
