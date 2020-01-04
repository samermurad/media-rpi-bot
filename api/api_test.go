package api

import (
	"encoding/json"
	"fmt"
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
