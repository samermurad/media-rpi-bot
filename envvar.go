package main

import (
	"fmt"
	"os"

	"samermurad.com/piBot/api"
)

type EnvVarCommand struct{}

func (envCmd *EnvVarCommand) Exec(data interface{}) error {
	if data, ok := data.(CmdExecData); ok {
		if len(data.Cmd.Args) > 0 {
			key := data.Cmd.Args[0]
			arg := os.Getenv(key)
			ch := make(chan *api.ApiResponse)
			msg := api.TelegramOutgoingMessage{
				ChatId:  data.Message.Chat.Id,
				Message: "Here is you arg:\nKey: " + key + "\nValue: " + arg,
			}
			go api.SendMessage(msg, ch)
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
