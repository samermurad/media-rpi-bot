package chatmachine

import "samermurad.com/piBot/telegram/models"

type ChatCache interface {
	SetTextMessage(string)
	GetTextMessage() string

	SetMessage(*models.BotMessage)
	GetMessage() *models.BotMessage

	Set(string, interface{})
	Get(string) interface{}

	ChatId() int64
}

type chatCache struct {
	dict    map[string]interface{}
	tmMsg   *models.BotMessage
	textMsg string
	chatId  int64
}

func (cc *chatCache) SetTextMessage(text string) {
	cc.textMsg = text
}

func (cc *chatCache) GetTextMessage() string {
	return cc.textMsg
}

func (cc *chatCache) SetMessage(msg *models.BotMessage) {
	cc.tmMsg = msg
}

func (cc *chatCache) GetMessage() *models.BotMessage {
	return cc.tmMsg
}

func (cc *chatCache) Set(key string, value interface{}) {
	if cc.dict == nil {
		cc.dict = make(map[string]interface{})
	}
	cc.dict[key] = value
}

func (cc *chatCache) Get(key string) interface{} {
	if cc.dict == nil {
		return nil
	}
	return cc.dict[key]
}

func (cc *chatCache) ChatId() int64 {
	return cc.chatId
}

func NewChatCache(chatId int64, initialMessage string) ChatCache {
	return &chatCache{
		chatId:  chatId,
		dict:    make(map[string]interface{}),
		textMsg: initialMessage,
	}
}
