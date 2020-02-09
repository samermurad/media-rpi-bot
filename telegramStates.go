package main

import (
	"fmt"
	"reflect"

	"www.samermurad.com/piBot/smachine"
)

type TelegramMessageState struct {
	count int
}

func (t *TelegramMessageState) OnEnter(prevState smachine.BaseState) {}
func (t *TelegramMessageState) OnExit(prevState smachine.BaseState)  {}

func (t *TelegramMessageState) Handle() reflect.Type {
	tmDebug(fmt.Sprintf("TelegramMessageState %v", t.count))
	t.count += 1
	return nil
}
