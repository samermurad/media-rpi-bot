package main

import (
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	os.Setenv("BOT_TOKEN", "STUB")
	// machine := smachine.NewStateMachine(map[reflect.Type]smachine.BaseState{
	// 	reflect.TypeOf(TelegramMessageState{}): &TelegramMessageState{},
	// })
	// machine.Start(1 * time.Second)
	// <-time.After(1 * time.Second)
	// machine.Stop()
}
