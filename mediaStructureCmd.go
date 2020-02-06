package main

import (
	"fmt"

	"samermurad.com/piBot/telegram"
	"samermurad.com/piBot/util"

	"samermurad.com/mediafilestructure/fs"
	"samermurad.com/piBot/telegram/models"
)

type MediaStructureCmd struct {
	SrcFolder  string
	DestFolder string
}

func (cmd *MediaStructureCmd) fileToMessageText(fileT fs.FileTransfer, editId int) string {
	return fmt.Sprintf(`
		*\=\=\=\=\= Title: %v \=\=\=\=*
		*ID: %v*
		*Current Location:*
		`+"`%v`"+`
		*New Location:*
		`+"`%v`"+`
	`,
		fileT.Title,
		editId,
		fileT.OldDir,
		fileT.NewDir,
	)
}
func (cmd *MediaStructureCmd) sendTransferForApproval(filesT []fs.FileTransfer, orginMsg *models.Message) {
	str := ""
	for i, f := range filesT {
		if i == 0 {
			str = cmd.fileToMessageText(f, i)
		} else {
			str += "\n" + cmd.fileToMessageText(f, i)
		}
	}
	ch := make(chan *models.Message)
	go telegram.SendMessage(
		models.BotMessage{
			Text:      str,
			ChatId:    orginMsg.Chat.Id,
			ParseMode: models.MarkdownV2,
		},
		ch,
	)
	msg := <-ch
	fmt.Println(msg)
}

func (cmd *MediaStructureCmd) Exec(data interface{}) error {
	if update, ok := data.(util.CmdExecData); ok {
		files := fs.WalkSource(cmd.SrcFolder, cmd.DestFolder)
		if len(files) == 0 {
			util.SendMessageAwait("No New Media", update.Message)
			return fmt.Errorf("No New Media")
		}
		fT := fs.ReadyFileTransfers(files, cmd.DestFolder)
		cmd.sendTransferForApproval(fT, update.Message)
	}
	return fmt.Errorf("Failed to parse command")
}

func (cmd *MediaStructureCmd) Args() map[string]interface{} {
	return map[string]interface{}{
		"src":  cmd.SrcFolder,
		"dest": cmd.DestFolder,
	}
}
