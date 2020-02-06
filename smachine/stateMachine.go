package smachine

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

type StateMachine interface {
	SetStates(states map[reflect.Type]BaseState)
	States() map[reflect.Type]BaseState

	Tick()

	Start() bool
	SetActive(bool)
	Stop()
}

type stateMachine struct {
	currentState   BaseState
	states         map[reflect.Type]BaseState
	isActive       bool
	runWG          sync.WaitGroup
	runCh          chan int
	updateInterval time.Duration
}

func (m *stateMachine) initState() {
	for _, v := range m.states {
		m.currentState = v
		break
	}
}

func (m *stateMachine) loop() {
	for {
		m.runWG.Add(1)
		select {
		case <-m.runCh:
			fmt.Println("State Mashine Stopped")
			m.runWG.Done()
			return
		default:
			fmt.Println("Starting Tick")
			if m.isActive {
				m.Tick()
				fmt.Println("Finished Tick, Waiting for", m.updateInterval)
			} else {
				fmt.Println("StateMaschine is paused", m.updateInterval)
			}
			<-time.After(m.updateInterval)
			m.runWG.Done()
		}
	}
}
func (m *stateMachine) SetStates(states map[reflect.Type]BaseState) {
	m.states = states
}

func (m *stateMachine) State() map[reflect.Type]BaseState {
	return m.states
}

func (m *stateMachine) Tick() {
	if m.currentState == nil {
		m.initState()
	}
	if typ := m.currentState.OnTick(); typ != nil && typ != reflect.TypeOf(m.currentState) {
		if newState := m.states[typ]; newState != nil {
			m.currentState.OnExit(newState)
			newState.OnEnter(m.currentState)
			m.currentState = newState
		}
	}
}

func (m *stateMachine) Start(updateInterval time.Duration) bool {
	if m.runCh == nil {
		m.runWG.Wait()
		fmt.Println("Start")
		m.updateInterval = updateInterval
		m.isActive = true
		m.runCh = make(chan int)
		go m.loop()
		return true
	} else {
		fmt.Println("Already Running")
		return false
	}
}
func (m *stateMachine) SetActive(isActive bool) {
	m.runWG.Wait()
	fmt.Println("Pause")
	m.isActive = true
}

func (m *stateMachine) Stop() {
	m.runWG.Wait()
	fmt.Println("Stop")
	m.isActive = false
	m.currentState = nil
	m.runCh <- 0
	m.runCh = nil
}

func NewStateMachine(states map[reflect.Type]BaseState) *stateMachine {
	m := new(stateMachine)
	m.runWG = sync.WaitGroup{}
	m.SetStates(states)
	return m
}
