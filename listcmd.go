package main

import (
	"fmt"
	"io/ioutil"
)

type ListCommand struct{}

func (lsCmd *ListCommand) Exec(data interface{}) error {
	if update, ok := data.(CmdExecData); ok {
		if len(update.Cmd.Args) > 0 {
			path := update.Cmd.Args[0]
			files, err := ioutil.ReadDir(path)
			if err != nil {
				SendMessageAwait(err.Error(), update.Message)
				return err
			} else {
				filesStr := "Here are the files in " + path + ":\n"
				for _, file := range files {
					dirExt, _ := (Ternary(file.IsDir(), `/`, "")).(string)
					filesStr += "\n" + file.Name() + dirExt
				}
				SendMessageAwait(filesStr, update.Message)
				return nil
			}
		} else {
			SendMessageAwait("Not Enough Vars", update.Message)
			return fmt.Errorf("Not Enough Vars")
		}
	} else {
		return fmt.Errorf("Failed to Parse Command")
	}
}

func (ls *ListCommand) Args() map[string]interface{} {
	panic("Not Implemented")
}
