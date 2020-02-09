package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"www.samermurad.com/piBot/spinner"
)

func Test1(t *testing.T) {
	os.Setenv("BOT_TOKEN", "STUB")
	sp := spinner.NewTmSpinner(68386493, "Timer Spinner")
	fmt.Println("WTF")
	chaa := make(chan bool)
	result := true
	for i := 0; i < 11; i++ {
		go sp.Progress(1, chaa)
		result = <-chaa
		<-time.After(2 * time.Second)
		if !result {
			t.Error("Bummer")
		}
	}
}
