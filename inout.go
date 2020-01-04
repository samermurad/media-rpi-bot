package main

import (
	"encoding/json"
	"fmt"
	"time"

	"samermurad.com/piBot/api"
	"samermurad.com/piBot/cmds"
	"samermurad.com/piBot/config"
	"samermurad.com/piBot/timeutils"
)

func attemptListeningToCmd(ch chan TrgmRes) {
	cb := make(chan (*api.ApiResponse))
	go api.GetUpdates(15*time.Second, config.CHAT_OFFSET(), cb)
	data := <-cb
	var obj api.TelegramGetUpdatesResponse
	err := json.Unmarshal(data.RawBody, &obj)

	if err != nil {
		ch <- TrgmRes{
			Res:   nil,
			Error: err,
		}
	} else {
		ch <- TrgmRes{
			Res:   &obj,
			Error: nil,
		}
	}
}


func Listener(dispatch chan<- *api.TelegramUpdate) {
	updateRes := make(chan TrgmRes)
	for {
		boom := time.After(500 * time.Millisecond)
		fmt.Println("Wait...")
		<-boom
		fmt.Println("Go!")
		go attemptListeningToCmd(updateRes)
		data := <-updateRes
		if data.Error != nil {
			fmt.Println("Error getting updates", data.Error)
		} else {
			if data.Res.Ok && len(data.Res.Result) > 0 {
				first := data.Res.Result[0]
				config.SET_CHAT_OFFSET(first.UpdateId + 1)
				dispatch <- &first
			}
		}
	}
}

func Handler(cmdMapping map[string]cmds.Command, source <-chan *api.TelegramUpdate) {
	startTime := timeutils.Seconds()
	for {
		fmt.Println("Handler before update := <-source")
		update := <-source
		fmt.Println("Handler after update := <-source")
		if startTime < update.Message.Date {
			if !IsChatAuthorized(update.Message.Chat.Id) {
				str := FetchRandomEvilInsult()
				tmDebug(str)
			} else {
				tmCmd, err := ParseCmdFromMessage(&update.Message)
				if err != nil {
					tmDebug("Bummer")
				} else if tmCmd == nil {
					tmDebug("No Cmd for u!")
				} else if cmd := cmdMapping[tmCmd.Key]; cmd != nil {
					data := CmdExecData{
						Message: &update.Message,
						Cmd:     tmCmd,
					}
					cmd.Exec(data)
				} else {
					tmDebug("Command not mapped")
				}
			}
		} else {
			fmt.Println("Ignoring old messages")
		}
	}
}
