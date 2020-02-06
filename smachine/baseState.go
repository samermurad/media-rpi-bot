package smachine

import (
	"fmt"
	"reflect"
)

type BaseState interface {
	OnTick() reflect.Type
	OnEnter(nextState BaseState)
	OnExit(prevState BaseState)
}

type printState struct{}

func (p *printState) OnEnter(prevState BaseState) {
	fmt.Println("Entering PrintState")
}

func (p *printState) OnExit(nextState BaseState) {
	fmt.Println("Exit PrintState")
}

func (p *printState) OnTick() reflect.Type {
	fmt.Println("Handle PrintState")
	return nil
}
