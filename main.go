package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	apiv2 "samermurad.com/piBot/api/v2"
	"samermurad.com/piBot/util"

	"samermurad.com/piBot/telegram"
	"samermurad.com/piBot/telegram/models"

	"samermurad.com/piBot/api"
	"samermurad.com/piBot/cmds"
	"samermurad.com/piBot/cntx"
)

type TrgmRes struct {
	Res   []*models.Update
	Error error
}

var updates chan TrgmRes
var uWg sync.WaitGroup
var quitWg = sync.WaitGroup{}
var startTime = time.Now()

var aContext cntx.ActionContext = cntx.ActionContext{}
var lastPinnedChatContext map[int64]string = make(map[int64]string)

var cmdMapping = map[string]cmds.Command{
	"/quit": &cmds.QuitCommand{
		Prepare: finishup,
		Token:   &quitWg,
	},
	"/env": &EnvVarCommand{},
	"/ls":  &ListCommand{},
	"/media_structure": &MediaStructureCmd{
		SrcFolder:  "/Users/samermurad/Movies/NewMedia",
		DestFolder: "/Users/samermurad/Movies/Dummy",
	},
}

func tmDebug(text string) {
	tgm := models.BotMessage{
		ChatId: 68386493,
		Text:   text,
	}
	channel := make(chan *models.Message)
	go telegram.SendMessage(tgm, channel)
	fmt.Println(text)
}

func getCntxKey(key string, chatId int64) string {
	return fmt.Sprintf("%v:%v", chatId, key)
}

func getEvilInsult() string {
	str := `U tryin to ask a bot how it's doing? LAMEEE`
	response := make(chan *apiv2.ResponseChannel)
	apiv2.NewBuilder("https://evilinsult.com/generate_insult.php?lang=en&type=json").
		Get().Build().
		Run(response)
	data := <-response
	if data.Err == nil {
		mp := make(map[string]interface{})
		if err := json.Unmarshal(data.Res.Body, &mp); err == nil && mp["insult"] != nil {
			s, _ := mp["insult"].(string)
			str = s
		}
	}
	return str
}

func getCmdFromMessage(msg *api.TelegramMssage) (*util.TMCommand, error) {
	cmd := &util.TMCommand{}
	msgText := msg.Text
	cmds := len(msg.Entities)
	if cmds == 0 {
		return nil, nil
	}
	if cmds > 1 {
		return nil, errors.New("Multi Commands not supported")
	}
	firstE := msg.Entities[0]
	cmd.Key = msgText[firstE.Offset:firstE.Length]
	cmd.Args = strings.Split(msgText[firstE.Length:], " ")
	return cmd, nil
}

func finishup() {
	tmDebug("Key byeeeee")
	<-time.After(2 * time.Second)
}

func main() {
	fmt.Println("Starting up..")
	updates := make(chan *models.Update)
	quitWg.Add(1)
	fmt.Println("Waiting for quit")
	go Listener(updates)
	go Handler(cmdMapping, updates)
	quitWg.Wait()
}
