package telegram

import (
	apiv2 "samermurad.com/piBot/api/v2"
	"samermurad.com/piBot/config"
	"samermurad.com/piBot/telegram/models"
)

const BASE_URL = "https://api.telegram.org/bot"
var apiBuilder = apiv2.NewBuilder()

func TlgrmUrl(path string) string {
	return BASE_URL + config.BOT_TOKEN() + path
}

func SendMessage(msg models.BotMessage, ch chan *models.ServerResponse) {
	// api.NewBuilder(BASE_URL).
	// Post().Build()
	apiv2.NewBuilder(TlgrmUrl("/sendMessage")).
		Post().
		MarshalBody(msg).
		

}
