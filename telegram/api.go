package telegram

import (
	"encoding/json"
	"fmt"
	"time"

	apiv2 "samermurad.com/piBot/api/v2"
	"samermurad.com/piBot/config"
	"samermurad.com/piBot/telegram/models"
)

const BASE_URL = "https://api.telegram.org/bot"

func apiBuilder() apiv2.RequestBuilder {
	return apiv2.
		NewBuilder(BASE_URL+config.BOT_TOKEN()).
		AddHeader("Content-Type", "application/json")
}

func run(req apiv2.Request, res models.Resultable, ch chan<- bool) {
	middle := make(chan *apiv2.ResponseChannel)
	go req.Run(middle)
	data := <-middle
	if data.Err != nil {
		ch <- false
	} else {
		// ok := new(models.OkResultCheck)
		if err := json.Unmarshal(data.Res.Body, res); err != nil && !res.IsOk() {
			fmt.Println(err)
			ch <- false
		} else {
			ch <- true
		}
	}
}
func SendMessage(msg models.BotMessage, ch chan<- *models.Message) {
	type wrapper struct {
		models.OkResultCheck
		Message models.Message `json:"result"`
	}
	res := new(wrapper)
	mid := make(chan bool)
	go run(
		apiBuilder().
			AppendUrl("/sendMessage").
			Post().
			MarshalBody(msg).
			Build(),
		res,
		mid,
	)
	ok := <-mid
	if ok {
		ch <- &res.Message
	} else {
		ch <- nil
	}
}

func EditMessageText(msg models.BotMessage, ch chan<- *models.Message) {
	res := new(models.MessageResult)
	mid := make(chan bool)
	go run(
		apiBuilder().
			AppendUrl("/editMessageText").
			Post().
			MarshalBody(msg).
			Build(),
		res,
		mid,
	)
	ok := <-mid
	if ok {
		ch <- &res.Message
	} else {
		ch <- nil
	}
}

func GetUpdates(updateOffset int64, timeout time.Duration, ch chan<- []*models.Update) {
	mp := make(map[string]interface{})
	mp["timeout"] = timeout
	mp["offset"] = updateOffset
	mid := make(chan bool)
	type wrapper struct {
		models.OkResultCheck
		List []*models.Update `json:"result"`
	}
	res := new(wrapper)
	go run(
		apiBuilder().
			AppendUrl("/getUpdates").
			Post().
			SetTimeout(timeout).
			MarshalBody(mp).
			Build(),
		res,
		mid,
	)
	ok := <-mid
	if !ok {
		ch <- nil
	} else {
		ch <- res.List
	}
}
