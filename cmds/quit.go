package cmds

import "sync"

type QuitCommand struct {
	Prepare func()
	Token   *sync.WaitGroup
}

func (qCmd *QuitCommand) Args() map[string]interface{} {
	return nil
}

func (qCmd *QuitCommand) Exec(data interface{}) error {
	qCmd.Prepare()
	qCmd.Token.Done()
	return nil
}
