package cmd

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/ssh/monitor"
)

type ObserveHandler struct {
	id uint64
	c  *gin.Context
}

func NewObserveHandler(c *gin.Context, id uint64) (_self Handler) {
	_self = &ObserveHandler{
		id: id,
		c:  c,
	}
	monitor.AddTerm(id)
	return
}

func (_self *ObserveHandler) GetIndex() int {
	return DefaultCmdHandler.GetIndex() + 2
}

func (_self *ObserveHandler) GetId() uint64 {
	return _self.id
}

func (_self *ObserveHandler) PassToClient(r rune) {
	monitor.Subject.NotifyUpdateObservers(_self.GetId(), r)
	return
}

func (_self *ObserveHandler) PassToServer(r rune) bool {
	return true
}

func (_self *ObserveHandler) Close() {
	monitor.RemoveTerm(_self.GetId())
}
