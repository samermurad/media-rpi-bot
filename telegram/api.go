package telegram

import (
	"encoding/json"
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

func run(req apiv2.Request, ch chan<- *models.ServerResponse) {
	middle := make(chan *apiv2.ResponseChannel)
	go req.Run(middle)
	data := <-middle
	if data.Err != nil {
		ch <- &models.ServerResponse{
			Ok:     false,
			Result: nil,
		}
	} else {
		res := new(models.ServerResponse)
		if err := json.Unmarshal(data.Res.Body, res); err != nil {
			ch <- &models.ServerResponse{
				Ok:     false,
				Result: nil,
			}
		} else {
			ch <- res
		}
	}
}
func SendMessage(msg models.BotMessage, ch chan<- *models.ServerResponse) {
	run(
		apiBuilder().
			AppendUrl("/sendMessage").
			Post().
			MarshalBody(msg).
			Build(),
		ch,
	)
}

func EditMessageText(msg models.BotMessage, ch chan<- *models.ServerResponse) {
	run(
		apiBuilder().
			AppendUrl("/editMessageText").
			Post().
			MarshalBody(msg).
			Build(),
		ch,
	)
}

func GetUpdates(updateOffset int64, timeout time.Duration, ch chan<- *models.ServerResponse) {
	mp := make(map[string]interface{})
	mp["timeout"] = timeout
	mp["offset"] = updateOffset
	run(
		apiBuilder().
			AppendUrl("/getUpdates").
			Post().
			SetTimeout(timeout).
			MarshalBody(mp).
			Build(),
		ch,
	)
}
