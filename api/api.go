package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"www.samermurad.com/piBot/config"
)

var BASE_URL string = "https://api.telegram.org/bot"

func apiRequest(req ApiRequest, ch chan *ApiResponse) {
	var url string
	httpLen := len("http")
	if len(req.Path) > httpLen && req.Path[0:httpLen] == "http" {
		url = req.Path // treat path as absolute URL
	} else {
		url = BASE_URL + config.BOT_TOKEN() + string(req.Path)
	}

	_req, err := http.NewRequest(req.Method, url, bytes.NewBuffer(req.Body))
	if err != nil {
		ch <- &ApiResponse{Error: ApiResponseError(err.Error())}
		return
	}
	req.SetHeaders(&_req.Header)
	fmt.Printf("Executing: %v %v\n%v\n", url, req.Method, string(req.Body))
	client := &http.Client{}
	if req.Timeout != 0 {
		client.Timeout = req.Timeout
	}
	response, err := client.Do(_req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		code := -1
		if response != nil && response.StatusCode != 0 {
			code = response.StatusCode
		}
		ch <- &ApiResponse{Error: ApiResponseError(err.Error()), StatusCode: code}
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		mp := make(map[string]interface{})
		err := json.Unmarshal(data, &mp)
		if err != nil {
			ch <- &ApiResponse{Error: ApiResponseError(err.Error()), StatusCode: response.StatusCode}
		} else {
			ch <- &ApiResponse{Body: mp, StatusCode: response.StatusCode, RawBody: data}
		}
	}
}

func SendMessage(msg TelegramOutgoingMessage, ch chan *ApiResponse) {
	req := NewApiReq()
	req.Path = "/sendMessage"
	req.Method = "POST"
	i, _ := json.Marshal(msg)
	req.Body = i
	go apiRequest(req, ch)
}

func EditMessageText(msg EditMessage, ch chan *ApiResponse) {
	req := NewApiReq()
	req.Path = "/editMessageText"
	req.Method = "POST"
	i, _ := json.Marshal(msg)
	req.Body = i
	go apiRequest(req, ch)
}
func GetUpdates(timeout time.Duration, offset int64, ch chan *ApiResponse) {
	req := NewApiReq()
	req.Path = "/getUpdates"
	req.Method = "POST"
	req.Timeout = timeout
	mp := make(map[string]interface{})
	mp["timeout"] = timeout
	mp["offset"] = offset
	_b, _ := json.Marshal(mp)
	req.Body = _b
	go apiRequest(req, ch)
}

func SendRequest(req ApiRequest, ch chan *ApiResponse) {
	go apiRequest(req, ch)
}
