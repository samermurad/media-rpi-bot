package models

type UpdateResults struct {
	OkResultCheck
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int64   `json:"update_id"`
	Message  Message `json:"message"`
}
