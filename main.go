package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"samermurad.com/piBot/api"
	"samermurad.com/piBot/cmds"
	"samermurad.com/piBot/cntx"
	"samermurad.com/piBot/config"
)

type TrgmRes struct {
	Res   *api.TelegramGetUpdatesResponse
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
}

func waitUpdates() {
	cb := make(chan (*api.ApiResponse))
	go api.GetUpdates(15*time.Second, config.CHAT_OFFSET(), cb)
	data := <-cb
	var obj api.TelegramGetUpdatesResponse
	err := json.Unmarshal(data.RawBody, &obj)
	uWg.Done()
	if err != nil {
		updates <- TrgmRes{
			Res:   nil,
			Error: err,
		}
	} else {
		updates <- TrgmRes{
			Res:   &obj,
			Error: nil,
		}
	}
}

func tmDebug(text string) {
	tgm := api.TelegramOutgoingMessage{
		ChatId:  68386493,
		Message: text,
	}
	channel := make(chan *api.ApiResponse)
	go api.SendMessage(tgm, channel)
	fmt.Println(text)
}

func getCntxKey(key string, chatId int64) string {
	return fmt.Sprintf("%v:%v", chatId, key)
}

func getEvilInsult() string {
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

func getCmdFromMessage(msg *api.TelegramMssage) (*TMCommand, error) {
	cmd := &TMCommand{}
	msgText := msg.Text
	cmds := len(msg.Entities)
	if cmds == 0 {
		return nil, nil
	}
	if cmds > 1 {
		return nil, api.ApiResponseError("Multi Commands not supported")
	}
	firstE := msg.Entities[0]
	cmd.Key = msgText[firstE.Offset:firstE.Length]
	cmd.Args = strings.Split(msgText[firstE.Length:], " ")
	return cmd, nil
}

func cmdAction(cmd *TMCommand, chatId int64) (followUpMessage string) {
	var aFn cntx.ActionFunction
	cntx := ""
	defer func() {
		if cntx != "" {
			aContext.Set(getCntxKey(cntx, chatId), &aFn)
			lastPinnedChatContext[chatId] = cntx
		}
	}()
	switch cmd.Key {
	case "/quit":
		cntx = "/quit"
		aFn = func(msg *api.TelegramMssage) *api.ApiResponse {
			ch := make(chan *api.ApiResponse)
			m := api.TelegramOutgoingMessage{}
			if config.APPROVAL_REG.Match([]byte(msg.Text)) {
				ch := make(chan *api.ApiResponse)
				m.Message = "Closing..."
				api.SendMessage(m, ch)
				res := <-ch
				return res
			} else {
				m.Message = "Cancelling.."
				api.SendMessage(m, ch)
				return <-ch
			}
		}
		return "Are you sure? (y/n)"
	default:
		break
	}
	return fmt.Sprintf("Command: %v Is Not Supported", cmd.Key)
}

func isChatAuthorized(cId int64) bool {
	for _, v := range config.ALLOWED_CHATS_IDS() {
		if v == cId {
			return true
		}
	}
	return false
}
func parseUpdateAction(update *api.TelegramUpdate, channel chan *api.ApiResponse) {
	config.SET_CHAT_OFFSET(update.UpdateId + 1)
	cmd, err := getCmdFromMessage(&update.Message)
	resMsg := api.TelegramOutgoingMessage{
		ChatId:  update.Message.Chat.Id,
		Message: "Unauthorized",
	}
	var followupFn func(msg *api.TelegramMssage) *api.ApiResponse
	defer func() {
		if followupFn != nil {
			go func() {
				res := followupFn(&update.Message)
				channel <- res
			}()
		} else {
			go api.SendMessage(resMsg, channel)
		}
	}()

	if !isChatAuthorized(update.Message.Chat.Id) {
		return
	}
	if lastCntxKey := lastPinnedChatContext[update.Message.Chat.Id]; lastCntxKey != "" {
		key := getCntxKey(lastCntxKey, update.Message.Chat.Id)
		followupFn = *aContext.Get(key)
		return
	}
	if err != nil {
		resMsg.Message = err.Error()
	} else if cmd != nil {
		resMsg.Message = cmdAction(cmd, update.Message.Chat.Id)
	} else {
		diss := getEvilInsult()
		resMsg.Message = diss
	}
}

func finishup() {
	tmDebug("Key byeeeee")
	<-time.After(2 * time.Second)
}

func main() {
	fmt.Println("Starting up..")
	updates := make(chan *api.TelegramUpdate)
	quitWg.Add(1)
	fmt.Println("Waiting for quit")
	go Listener(updates)
	go Handler(cmdMapping, updates)
	quitWg.Wait()
}
