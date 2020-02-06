package chatmachine

import (
	"io/ioutil"

	"samermurad.com/piBot/util"
)

type LsState struct {
}

func (state *LsState) Exec(data util.CmdExecData, cache ChatCache) ChatState {
	if len(data.Cmd.Args) > 0 {
		path := data.Cmd.Args[0]
		files, err := ioutil.ReadDir(path)
		if err != nil {
			cache.SetTextMessage(err.Error())
			return nil
		} else {
			filesStr := "Here are the files in " + path + ":\n"
			for _, file := range files {
				dirExt, _ := (util.Ternary(file.IsDir(), `/`, "")).(string)
				filesStr += "\n" + file.Name() + dirExt
			}
			cache.SetTextMessage(filesStr)
			return nil
		}
	} else {
		cache.SetTextMessage("I require some Args")
	}
	return nil
}
