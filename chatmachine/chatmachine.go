package chatmachine

import (
	"fmt"

	"samermurad.com/piBot/telegram"

	"samermurad.com/piBot/util"

	"samermurad.com/piBot/config"
	"samermurad.com/piBot/telegram/models"
)

const (
	CACHE_KEY_DONE_MESSAGE     = "CHAT_MACHINE_SHUTDOWN_MESSAGE"
	CACHE_KEY_TELEGRAM_MESSAGE = "CACHE_KEY_TELEGRAM_MESSAGE"
)

func IsChatAuthorized(cId int64) bool {
	for _, v := range config.ALLOWED_CHATS_IDS() {
		if v == cId {
			return true
		}
	}
	return false
}

type ChatState interface {
	Exec(update util.CmdExecData, cache ChatCache) ChatState
}

type ChatMachine struct {
	ChatId         int64
	Cache          ChatCache
	CurrentState   ChatState
	States         map[string]ChatState
	UpdateChannel  <-chan interface{}
	DoneChannel    chan<- bool
	ShutDownMesage string
}

type TelegramUpdateJob struct {
	Done           bool
	ChatId         int64
	UpdateChannel  <-chan interface{}
	DoneChannel    chan<- bool
	ShutDownMesage string
}

func (chm *ChatMachine) Start(updates <-chan interface{}, done chan<- bool) {
	for u := range updates {
		fmt.Println("Got Update")
		tu := u.(*models.Update)
		done <- chm.Run(tu)
	}
	fmt.Println("Finished")
}

func (chm *ChatMachine) ShutDown() {
	fmt.Println("Shutting Down")
	if tmMsg := chm.Cache.GetMessage(); tmMsg != nil {
		ch := make(chan *models.Message)
		go telegram.SendMessage(*tmMsg, ch)
		<-ch
	} else {
		msgText := chm.Cache.GetTextMessage()
		util.SendQuickTelegramMessage(chm.ChatId, msgText)
	}
}

func (chm *ChatMachine) Run(update *models.Update) bool {
	if chm.Cache.GetTextMessage() == "" {
		chm.Cache.SetTextMessage("Action Timedout")
	}
	if !IsChatAuthorized(update.Message.Chat.Id) {
		chm.Cache.SetTextMessage("I Dont know you, go away")
		return true
	}

	tmCmd, _ := util.ParseCmdFromMessage(&update.Message)
	data := util.CmdExecData{Message: &update.Message}

	if tmCmd == nil {
		if chm.CurrentState == nil {
			evil := util.FetchRandomEvilInsult()
			chm.Cache.SetTextMessage(evil)
			return true
		} else {
			if nextState := chm.CurrentState.Exec(data, chm.Cache); nextState != nil {
				chm.CurrentState = nextState
				return false
			} else {
				return true
			}
		}
	} else {
		data.Cmd = tmCmd
		if chm.CurrentState == nil {
			chm.CurrentState = chm.States[tmCmd.Key]
		}
		if chm.CurrentState != nil {
			if nextState := chm.CurrentState.Exec(data, chm.Cache); nextState != nil {
				chm.CurrentState = nextState
				return false
			} else {
				return true
			}
		} else {
			chm.Cache.SetTextMessage("Command not mapped, and im not just gonna start doing magic")
			return true
		}
	}
}
