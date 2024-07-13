package cmd

import (
	"go-protector/server/internal/ssh/notify"
	"go-protector/server/internal/ws"
)

type NotifyCliHandler struct {
	id uint64
}

func NewNotifyHandler(impl ws.IWsWriteCloser, id uint64) Handler {
	notify.RegisterWriter(id, impl)
	return &NotifyCliHandler{
		id: id,
	}
}

func (_self *NotifyCliHandler) GetIndex() int {
	return DefaultCmdHandler.GetIndex() + 3
}

func (_self *NotifyCliHandler) GetId() uint64 {
	return _self.id
}

func (_self *NotifyCliHandler) PassToClient(r rune) {
	return
}

func (_self *NotifyCliHandler) PassToServer(r rune) bool {
	return true
}

func (_self *NotifyCliHandler) Close() {
	notify.UnRegisterWriter(_self.GetId())
}
