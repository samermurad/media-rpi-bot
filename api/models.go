package api

import (
	"net/http"
	"time"
)

type ApiResponseError string

func (aErr ApiResponseError) Error() string {
	return string(aErr)
}

type ApiResponse struct {
	StatusCode int
	Error      ApiResponseError
	Body       map[string]interface{}
	RawBody    []byte
}

type ApiRequestHeaders map[string]string

type ApiRequest struct {
	Path    string
	Headers ApiRequestHeaders
	Method  string
	Body    []byte
	Timeout time.Duration
}

type TelegramOutgoingMessage struct {
	ChatId  int64  `json:"chat_id"`
	Message string `json:"text"`
}

type EditMessage struct {
	ChatId    int64  `json:"chat_id"`
	Text      string `json:"text"`
	MessageId int64  `json:"message_id"`
}

func (req ApiRequest) SetHeaders(hdrs *http.Header) {
	for k, v := range req.Headers {
		hdrs.Set(k, v)
	}
}

func NewApiReq() ApiRequest {
	hdrs := make(ApiRequestHeaders)
	hdrs["Content-Type"] = "application/json"
	return ApiRequest{
		Method:  "GET",
		Headers: hdrs,
	}
}
