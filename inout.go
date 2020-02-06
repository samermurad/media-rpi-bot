package main

import (
	"fmt"
	"time"

	"samermurad.com/piBot/chatmachine"
	"samermurad.com/piBot/telegram"
	"samermurad.com/piBot/telegram/models"
	"samermurad.com/piBot/timeutils"
	"samermurad.com/piBot/util"

	"samermurad.com/jobDispatcher/dispatch"
	"samermurad.com/piBot/cmds"
	"samermurad.com/piBot/config"
)

func attemptListeningToCmd(ch chan TrgmRes) {
	cb := make(chan ([]*models.Update))
	go telegram.GetUpdates(config.CHAT_OFFSET(), 15*time.Second, cb)
	data := <-cb
	ch <- TrgmRes{
		Res:   data,
		Error: nil,
	}
}

func Listener(dispatch chan<- *models.Update) {
	updateRes := make(chan TrgmRes)
	for {
		boom := time.After(500 * time.Millisecond)
		fmt.Println("Wait...")
		<-boom
		fmt.Println("Go!")
		go attemptListeningToCmd(updateRes)
		data := <-updateRes
		if data.Error != nil {
			fmt.Println("Error getting updates", data.Error)
		} else {
			if data.Res != nil {
				first := data.Res[0]
				config.SET_CHAT_OFFSET(first.UpdateId + 1)
				dispatch <- first
			}
		}
	}
}

type TelegramUpdateJob struct {
	Done           bool
	ChatId         int64
	UpdateChannel  <-chan interface{}
	DoneChannel    chan<- bool
	ShutDownMesage string
}

func (tJ *TelegramUpdateJob) Start(updates <-chan interface{}, done chan<- bool) {
	for u := range updates {
		fmt.Println("Got Update")
		tu := u.(*models.Update)
		done <- tJ.Run(tu)
	}
}

func (tJ *TelegramUpdateJob) ShutDown() {
	fmt.Println("Shutting Down")
	tmDebug(tJ.ShutDownMesage)
}

func (tJ *TelegramUpdateJob) Run(update *models.Update) bool {
	if !util.IsChatAuthorized(update.Message.Chat.Id) {
		// 		str := FetchRandomEvilInsult()
		tJ.ShutDownMesage = "I Dont know you, go away"
		return true
	}

	tmCmd, err := util.ParseCmdFromMessage(&update.Message)
	if err != nil {
		tJ.ShutDownMesage = "Bummer"
		return true
	} else if tmCmd == nil {
		tJ.ShutDownMesage = "No Cmd for u!"
		return true
	} else if cmd := cmdMapping[tmCmd.Key]; cmd == nil {
		tJ.ShutDownMesage = "Command not mapped"
		return true
	} else {
		data := util.CmdExecData{
			Message: &update.Message,
			Cmd:     tmCmd,
		}
		cmd.Exec(data)
		return true
	}
}

func Handler(cmdMapping map[string]cmds.Command, source <-chan *models.Update) {
	startTime := timeutils.Seconds()
	var creator func(update interface{}) (dispatch.Job, time.Duration)
	creator = func(update interface{}) (dispatch.Job, time.Duration) {
		tmUpdate, ok := update.(*models.Update)
		if !ok {
			return nil, 0
		}
		if startTime < tmUpdate.Message.Date {
			return &chatmachine.ChatMachine{
				ChatId: tmUpdate.Message.Chat.Id,
				Cache:  chatmachine.NewChatCache(tmUpdate.Message.Chat.Id, "Action Timed Out"),
				States: map[string]chatmachine.ChatState{
					"/ls": &chatmachine.LsState{},
					"/media_structure": &chatmachine.OrganizeMedia{
						SrcFolder:  "/Users/samermurad/Movies/NewMedia",
						DestFolder: "/Users/samermurad/Movies/Dummy",
					},
					"/cmd": &chatmachine.RunSystemScript{
						Cmd:  "source $HOME/.bashrc && cpip",
						Args: nil,
					},
					"/sleep": &chatmachine.RunSystemScript{
						Cmd:  "sleep 2 && echo \"Slept for 2 secs\"",
						Args: nil,
					},
				},
			}, 10 * time.Second
		} else {
			fmt.Println("Ignoring old messages")
			return nil, 0
		}
		return nil, 0
	}
	dispatch := dispatch.JobDispatcher{JobCreator: creator}
	for {
		fmt.Println("Handler before update := <-source")
		update := <-source
		dispatch.Dispatch(update.Message.Chat.Id, update)
		// fmt.Println("Handler after update := <-source")
		// if startTime < update.Message.Date {
		// 	if !IsChatAuthorized(update.Message.Chat.Id) {
		// 		str := FetchRandomEvilInsult()
		// 		tmDebug(str)
		// 	} else {
		// 		tmCmd, err := ParseCmdFromMessage(&update.Message)
		// 		if err != nil {
		// 			tmDebug("Bummer")
		// 		} else if tmCmd == nil {
		// 			tmDebug("No Cmd for u!")
		// 		} else if cmd := cmdMapping[tmCmd.Key]; cmd != nil {
		// 			data := CmdExecData{
		// 				Message: &update.Message,
		// 				Cmd:     tmCmd,
		// 			}
		// 			cmd.Exec(data)
		// 		} else {
		// 			tmDebug("Command not mapped")
		// 		}
		// 	}
		// } else {
		// 	fmt.Println("Ignoring old messages")
		// }
	}
}
