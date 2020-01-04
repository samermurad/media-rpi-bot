package api

import "time"

/**
{
            "update_id": 65651909,
            "message": {
                "message_id": 11,
                "from": {
                    "id": 68386493,
                    "is_bot": false,
                    "first_name": "Samer",
                    "last_name": "Murad",
                    "username": "SamerMurad",
                    "language_code": "en"
                },
                "chat": {
                    "id": 68386493,
                    "first_name": "Samer",
                    "last_name": "Murad",
                    "username": "SamerMurad",
                    "type": "private"
                },
                "date": 1577447969,
                "text": "/later yeah",
                "entities": [
                    {
                        "offset": 0,
                        "length": 6,
                        "type": "bot_command"
                    }
                ]
            }
        }
*/

type TelegramFrom struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
	LangCode  string `json:"language_code"`
}

type TelegramChat struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type TelegramMessageEntity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type TelegramMssage struct {
	MessageId int64                   `json:"message_id"`
	From      TelegramFrom            `json:"from"`
	Chat      TelegramChat            `json:"chat"`
	Entities  []TelegramMessageEntity `json:"entities"`
	Text      string                  `json:"text"`
	Date      time.Duration           `json:"date"`
}

type TelegramUpdate struct {
	UpdateId int64          `json:"update_id"`
	Message  TelegramMssage `json:"message"`
}

type TelegramGetUpdatesResponse struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

type TelegramSenMessageResponse struct {
	Ok     bool           `json:"ok"`
	Result TelegramMssage `json:"result"`
}
