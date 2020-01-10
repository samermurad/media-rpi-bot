package main

import (
	"reflect"
	"testing"
	"time"

	"samermurad.com/piBot/smachine"
)

func TestBasic(t *testing.T) {
	machine := smachine.NewStateMachine(map[reflect.Type]smachine.BaseState{
		reflect.TypeOf(TelegramMessageState{}): &TelegramMessageState{},
	})
	machine.Start(1 * time.Second)
	<-time.After(4 * time.Second)
	// machine.Stop()
}
