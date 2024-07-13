package monitor

import (
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/ws"
)

type Observer struct {
	ObId  uint64
	SsoId uint64
	ws.IWsWriter
}

func (_self *Observer) GetObId() uint64 {
	return _self.ObId
}

func (_self *Observer) GetSsoId() uint64 {
	return _self.SsoId
}
func (_self *Observer) Update(str string) {
	_ = _self.Write(base.NewWsMsg(consts.MsgData, str))
	return
}
func (_self *Observer) Close() {
	_ = _self.Write(base.NewWsMsg(consts.MsgClose, "连接已关闭"))
	return
}
