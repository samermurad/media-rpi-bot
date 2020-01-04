package cntx

import (
	"sync"

	"samermurad.com/piBot/api"
)

type ActionFunction func(*api.TelegramMssage) *api.ApiResponse
type ActionContext struct {
	Mu   sync.Mutex
	cntx map[string]*ActionFunction
}

func (aCntx *ActionContext) Get(key string) *ActionFunction {
	defer aCntx.Mu.Unlock()
	aCntx.Mu.Lock()
	if aCntx.cntx[key] != nil {
		return aCntx.cntx[key]
	}
	return nil
}

func (aCntx *ActionContext) Set(key string, value *ActionFunction) {
	defer aCntx.Mu.Unlock()
	aCntx.Mu.Lock()
	if aCntx.cntx == nil {
		aCntx.cntx = make(map[string]*ActionFunction)
	}
	aCntx.cntx[key] = value
}
