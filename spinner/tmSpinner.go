package spinner

import (
	"fmt"

	"www.samermurad.com/piBot/telegram"
	"www.samermurad.com/piBot/telegram/models"
)

type Spinner interface {
	Prepare(channel chan bool)
	Progress(steps int, channel chan bool)
	Steps() int
	Finish(channel chan bool)
}

type spinner struct {
	description string
	steps       int
	totalSteps  int
	tmMsgId     int64
	tmChatId    int64
	didFinish   bool
}

func getBoolCh() chan bool {
	chBool := make(chan bool)
	return chBool
}
func getApiResCh() chan *models.Message {
	chBool := make(chan *models.Message)
	return chBool
}

func getMsg(text string, chId int64) *models.BotMessage {
	return &models.BotMessage{
		Text:   text,
		ChatId: chId,
	}
}

func (sp *spinner) Prepare(channel chan bool) {
	msg := getMsg("Preparing "+sp.description, sp.tmChatId)
	apiCh := getApiResCh()
	go telegram.SendMessage(*msg, apiCh)
	if data := <-apiCh; data == nil {
		channel <- false
		return
	} else {
		sp.tmMsgId = data.MessageId
		channel <- true
		return
	}
}

func (sp *spinner) getLoaderString() string {
	loadingStep := "---"
	readyStep := "==="
	loaderStr := ""
	fullStr := "Loading: " + sp.description + "\n\n"
	for i := 0; i < sp.totalSteps; i++ {
		if sp.steps < i {
			loaderStr += loadingStep
		} else {
			loaderStr += readyStep
		}
	}
	fullStr += loaderStr + " " + fmt.Sprintf("%v/%v", sp.steps+1, sp.totalSteps)
	return fullStr
}
func (sp *spinner) Finish(channel chan bool) {
	if sp.didFinish {
		channel <- false
		fmt.Println("Cannot finish an already finished spinner")
		return
	}
	sp.didFinish = true
	msg := getMsg("Done: "+sp.description, sp.tmChatId)
	ch := getApiResCh()
	go telegram.SendMessage(*msg, ch)
	<-ch
	channel <- true
}

func (sp *spinner) Steps() int {
	return sp.steps
}

func (sp *spinner) Progress(steps int, channel chan bool) {
	if sp.steps >= sp.totalSteps {
		sp.Finish(channel)
		return
	}
	str := sp.getLoaderString()
	msg := &models.BotMessage{
		ChatId:    sp.tmChatId,
		MessageId: sp.tmMsgId,
		Text:      str,
	}
	apiCh := getApiResCh()
	go telegram.EditMessageText(*msg, apiCh)

	if data := <-apiCh; data == nil {
		channel <- false
	} else {
		sp.steps += steps
		channel <- true
	}
}

func (sp *spinner) Invalidate(channel chan bool) {
}

func NewTmSpinner(chatId int64, itemDesc string) Spinner {
	sp := spinner{}
	sp.tmChatId = chatId
	sp.steps = 0
	sp.totalSteps = 10
	sp.description = itemDesc
	chBool := getBoolCh()
	go sp.Prepare(chBool)
	<-chBool
	return &sp
}
