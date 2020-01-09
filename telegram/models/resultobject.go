package models

type Resultable interface {
	IsOk() bool
}

type OkResultCheck struct {
	Ok bool `json:"ok"`
}

func (ok *OkResultCheck) IsOk() bool {
	return ok.Ok
}
