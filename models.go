package main

import (
	"fmt"

	"samermurad.com/piBot/api"
)

type TMCommand struct {
	Key  string
	Args []string
}

func (cmd *TMCommand) String() string {
	return fmt.Sprintf("Command: %v, Args: %v", cmd.Key, cmd.Args)
}

type CmdExecData struct {
	Message *api.TelegramMssage
	Cmd     *TMCommand
}
