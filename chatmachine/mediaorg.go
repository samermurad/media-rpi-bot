package chatmachine

import (
	"fmt"
	"time"

	"www.samermurad.com/mediafilestructure/fs"
	"www.samermurad.com/piBot/telegram"
	"www.samermurad.com/piBot/telegram/models"
	"www.samermurad.com/piBot/util"
)

type OrganizeMediaState int

const (
	_organizeMediaStateInitial       OrganizeMediaState = 0
	_organizeMediaStateAwaitApproval OrganizeMediaState = 1
)

type OrganizeMedia struct {
	SrcFolder  string
	DestFolder string
	state      OrganizeMediaState
	files      []fs.FileTransfer
}

func (cmd *OrganizeMedia) fileToMessageText(fileT fs.FileTransfer, editId int) string {
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
func (cmd *OrganizeMedia) sendTransferForApproval(filesT []fs.FileTransfer, orginMsg *models.Message) {
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

func (cmd *OrganizeMedia) unleashHell(msg *models.Message) {
	for i := range time.Tick(1 * time.Second) {
		util.SendMessageAwait(fmt.Sprint("This method is going to unleash hell", i), msg)
	}
}

func (cmd *OrganizeMedia) GetYesNo(msg *models.Message) models.BotMessage {
	btns := [][]models.KeyboardButton{
		[]models.KeyboardButton{
			models.KeyboardButton{Text: "Yes"},
			models.KeyboardButton{Text: "No"},
		},
	}
	return models.BotMessage{
		Text:   "Does this seem okay?",
		ChatId: msg.Chat.Id,
		ReplyMarkup: models.ReplyKeyboardMarkup{
			Keyboard:       btns,
			OnTimeKeyboard: true,
			Selective:      false,
		},
	}
}

func (cmd *OrganizeMedia) GetResetMessage(text string, msg *models.Message) *models.BotMessage {
	return &models.BotMessage{
		Text:   text,
		ChatId: msg.Chat.Id,
		ReplyMarkup: models.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      false,
		},
	}
}
func (cmd *OrganizeMedia) transferFiles(msg *models.Message) {
	fs.TransferFiles(cmd.files)
	util.SendMessageAwait("Done!", msg)
}
func (cmd *OrganizeMedia) Exec(data util.CmdExecData, cache ChatCache) ChatState {
	if cmd.state == _organizeMediaStateInitial {
		files := fs.WalkSource(cmd.SrcFolder, cmd.DestFolder)
		if len(files) == 0 {
			cache.SetTextMessage("No New Media")
			return nil
		}
		cmd.files = fs.ReadyFileTransfers(files, cmd.DestFolder)
		cmd.sendTransferForApproval(cmd.files, data.Message)
		util.SendBotMessageAwait(cmd.GetYesNo(data.Message))
		cmd.state = _organizeMediaStateAwaitApproval
		cache.SetMessage(cmd.GetResetMessage("You Waited too long to give an answer... reseting state", data.Message))
		return cmd
	} else if cmd.state == _organizeMediaStateAwaitApproval {
		if data.Message.Text == "Yes" {
			cache.SetMessage(cmd.GetResetMessage("You Pressed Yes!, transferring files...", data.Message))
			go cmd.transferFiles(data.Message)
			return nil
		} else if data.Message.Text == "No" {
			cache.SetMessage(cmd.GetResetMessage("Aborting", data.Message))
			return nil
		}
		util.SendMessageAwait("Inavlid input", data.Message)
		return cmd
	} else {
	}
	return nil
}
