package chatmachine

import (
	"fmt"
	"os/exec"
	"sync"

	"samermurad.com/piBot/telegram"

	"samermurad.com/piBot/telegram/models"
	"samermurad.com/piBot/util"
)

type RunSystemScript struct {
	Cmd  string
	Args []string
}

const shellToUse = "bash"

type tmBuffer struct {
	chatId int64
	btMsg  *models.BotMessage
}

var mu sync.Mutex = sync.Mutex{}

func (tmB *tmBuffer) getBotMsg() *models.BotMessage {
	if tmB.btMsg == nil {
		tmB.btMsg = &models.BotMessage{
			ChatId: tmB.chatId,
		}
	}
	return tmB.btMsg
}

func (tmB *tmBuffer) Update(txt string) *models.Message {
	msg := tmB.getBotMsg()
	msg.Text = txt
	ch := make(chan *models.Message)
	go telegram.EditMessageText(*msg, ch)
	return <-ch
}
func (tmB *tmBuffer) Write(data []byte) (int, error) {
	str := string(data)
	tmB.Update(str)
	return len(data), nil
}

func (script *RunSystemScript) shellout(command string, tmb *tmBuffer) error {
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	mu.Lock()
	cmd := exec.Command(shellToUse, "-c", command)
	cmd.Stdout = tmb
	cmd.Stderr = tmb
	err := cmd.Run()
	mu.Unlock()
	return err //stdout.String(), stderr.String()
}

func (script *RunSystemScript) kickStart(data util.CmdExecData) {
	buffer := &tmBuffer{
		chatId: data.Message.Chat.Id,
	}
	firstMsgChan := make(chan *models.Message)
	firstMsg := buffer.getBotMsg()
	firstMsg.Text = fmt.Sprint("Task:", script.Cmd, "Is Pending")
	go telegram.SendMessage(*firstMsg, firstMsgChan)
	res := <-firstMsgChan
	go script.shellout(script.Cmd, buffer)
	buffer.btMsg.MessageId = res.MessageId
}
func (script *RunSystemScript) Exec(data util.CmdExecData, cache ChatCache) ChatState {
	cache.SetTextMessage(fmt.Sprint("Task sent to background"))
	go script.kickStart(data)
	return nil
}
