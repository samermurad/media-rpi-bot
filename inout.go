package main

import (
	"fmt"
	"time"

	"samermurad.com/piBot/telegram"
	"samermurad.com/piBot/telegram/models"

	"samermurad.com/piBot/cmds"
	"samermurad.com/piBot/config"
	"samermurad.com/piBot/timeutils"
)

func attemptListeningToCmd(ch chan TrgmRes) {
	cb := make(chan ([]*models.Update))
	go telegram.GetUpdates(config.CHAT_OFFSET(), 15*time.Second, cb)
	data := <-cb
	ch <- TrgmRes{
		Res:   data,
		Error: nil,
	}
}

func Listener(dispatch chan<- *models.Update) {
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
			if data.Res != nil {
				first := data.Res[0]
				config.SET_CHAT_OFFSET(first.UpdateId + 1)
				dispatch <- first
			}
		}
	}
}

func Handler(cmdMapping map[string]cmds.Command, source <-chan *models.Update) {
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
