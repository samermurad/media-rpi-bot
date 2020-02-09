package main

import (
	"fmt"
	"os"

	"www.samermurad.com/piBot/telegram"
	tmM "www.samermurad.com/piBot/telegram/models"
	"www.samermurad.com/piBot/util"
)

type EnvVarCommand struct{}

func (envCmd *EnvVarCommand) Exec(data interface{}) error {
	if data, ok := data.(util.CmdExecData); ok {
		if len(data.Cmd.Args) > 0 {
			key := data.Cmd.Args[0]
			arg := os.Getenv(key)
			ch := make(chan *tmM.Message)
			msg := tmM.BotMessage{
				ChatId: data.Message.Chat.Id,
				Text:   "Here is you arg:\nKey: " + key + "\nValue: " + arg,
			}
			go telegram.SendMessage(msg, ch)
			<-ch
		}
	} else {
		return fmt.Errorf("Error, Cmd missing arg")
	}
	return nil
}

func (envCmd *EnvVarCommand) Args() map[string]interface{} {
	panic("Not Implemented")
}
