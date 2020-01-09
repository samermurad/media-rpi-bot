package main

import (
	"fmt"
	"strings"

	"samermurad.com/piBot/telegram"
	"samermurad.com/piBot/telegram/models"

	"samermurad.com/piBot/api"
	"samermurad.com/piBot/config"
)

type StrArr []string

func (arr StrArr) FilterEmpty() []string {
	strs := make([]string, 0)
	for _, v := range arr {
		if v != "" {
			strs = append(strs, v)
		}
	}
	return strs
}

func IsChatAuthorized(cId int64) bool {
	for _, v := range config.ALLOWED_CHATS_IDS() {
		if v == cId {
			return true
		}
	}
	return false
}
func ParseCmdFromMessage(msg *models.Message) (*TMCommand, error) {
	cmd := &TMCommand{}
	msgText := msg.Text
	cmds := len(msg.Entities)
	if cmds == 0 {
		return nil, fmt.Errorf("No Command Found")
	}
	if cmds > 1 {
		return nil, fmt.Errorf("Multi Commands not supported")
	}
	firstE := msg.Entities[0]
	cmd.Key = msgText[firstE.Offset:firstE.Length]
	cmd.Args = StrArr(strings.Split(msgText[firstE.Length:], " ")).FilterEmpty()
	return cmd, nil
}

func FetchRandomEvilInsult() string {
	str := `U tryin to ask a bot how it's doing? LAMEEE`
	insultApi := api.ApiRequest{
		Path:   "https://evilinsult.com/generate_insult.php?lang=en&type=json",
		Method: "GET",
	}
	response := make(chan *api.ApiResponse)
	api.SendRequest(insultApi, response)
	data := <-response
	if insult := data.Body["insult"]; insult != nil {
		s, _ := insult.(string)
		str = s
	}
	return str
}

func MakeContextKey(key string, chatId int64) string {
	return fmt.Sprintf("%v:%v", chatId, key)
}

func SendMessageAwait(text string, tmMsg *models.Message) chan *models.Message {
	msg := models.BotMessage{
		ChatId: tmMsg.Chat.Id,
		Text:   text,
	}
	ch := make(chan *models.Message)
	go telegram.SendMessage(msg, ch)
	return ch
}

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}
