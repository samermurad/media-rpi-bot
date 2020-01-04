package api

import (
	"sync"
	"testing"
	"time"
)

func TestWeird(t *testing.T) {
	s := sync.WaitGroup{}
	var wtc = func(delay time.Duration, ch chan int) {
		<-time.After(delay * time.Second)
		ch <- int(delay)
	}
	s.Add(2)
	go func() {
		o := make(chan int)
		println("O Created")
		go wtc(1, o)
		println("O Will Read")
		data := <-o
		println("O Got the ", data)
		println("O Closing")
		s.Done()
	}()
	go func() {
		t := make(chan int)
		println("T Created")
		wtc(1, t)
		println("T Will Read")
		data := <-t
		println("T Got the ", data)
		println("T Closing")
		s.Done()
	}()
	println("Starting")
	s.Wait()
}
