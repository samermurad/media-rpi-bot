package models

import "time"

type BotMessageParseMode string

const (
	HTML       BotMessageParseMode = "HTML"
	Markdown   BotMessageParseMode = "Markdown"
	MarkdownV2 BotMessageParseMode = "MarkdownV2"
)

type MessageResult struct {
	OkResultCheck
	Message Message `json:"result"`
}
type Message struct {
	MessageId int64         `json:"message_id"`
	From      From          `json:"from"`
	Chat      Chat          `json:"chat"`
	Entities  []TextEntity  `json:"entities"`
	Text      string        `json:"text"`
	Date      time.Duration `json:"date"`
}

type BotMessage struct {
	ChatId           int64               `json:"chat_id"`
	Text             string              `json:"text"`
	ParseMode        BotMessageParseMode `json:"parse_mode,omitempty"`
	MessageId        int64               `json:"message_id,omitempty"`
	ReplyMarkup      ReplyMarkup         `json:"reply_markup,omitempty"`
	ReplyToMessageId int64               `json:"reply_to_message_id,omitempty"`
}
