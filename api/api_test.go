package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"samermurad.com/piBot/config"
)

func TestSendMessage(t *testing.T) {
	chtId := int64(68386493)
	msg := TelegramOutgoingMessage{
		ChatId:  chtId,
		Message: "Testing Golang Telegram Bot integration",
	}
	kay := make(chan *ApiResponse)
	SendMessage(msg, kay)
	data := <-kay
	fmt.Println(data)
	if data.StatusCode != 200 {
		t.Errorf("data.StatusCode != 200")
	}
}

func TestGetUpdates(t *testing.T) {
	// chtId := int64(68386493)
	kay := make(chan *ApiResponse)
	GetUpdates(15*time.Second, config.CHAT_OFFSET(), kay)
	data := <-kay
	fmt.Println(data.Body)
	if data.StatusCode != 200 {
		t.Errorf("data.StatusCode != 200")
	} else {
		msgs := TelegramGetUpdatesResponse{}
		err := json.Unmarshal(data.RawBody, &msgs)
		if err != nil {
			t.Errorf(err.Error())
		} else {
			fmt.Println(msgs)
		}
	}
}
func TestWeird(t *testing.T) {
	s := sync.WaitGroup{}
	var wtc = func(delay time.Duration, ch chan int) {
		<-time.After(delay * time.Second)
		ch <- int(delay)
	}
	s.Add(2)
	go func() {
		o := make(chan int)
		println("O Created")
		go wtc(1, o)
		println("O Will Read")
		data := <-o
		println("O Got the ", data)
		println("O Closing")
		s.Done()
	}()
	go func() {
		t := make(chan int)
		println("T Created")
		wtc(1, t)
		println("T Will Read")
		data := <-t
		println("T Got the ", data)
		println("T Closing")
		s.Done()
	}()
	println("Starting")
	s.Wait()
}
