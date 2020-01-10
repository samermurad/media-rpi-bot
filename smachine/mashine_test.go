package smachine

import (
	"reflect"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	machine := NewStateMachine(map[reflect.Type]BaseState{
		reflect.TypeOf(printState{}): &printState{},
	})
	machine.Start(1 * time.Second)
	<-time.After(4 * time.Second)
}
