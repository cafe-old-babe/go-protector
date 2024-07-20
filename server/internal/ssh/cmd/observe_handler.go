package cmd

import (
	"bytes"
	"context"
	"go-protector/server/internal/ssh/monitor"
	"go-protector/server/internal/utils/async"
	"time"
	"unicode/utf8"
)

type ObserveHandler struct {
	id       uint64
	dataChan chan rune
	*time.Ticker
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewObserveHandler(c context.Context, id uint64) Handler {
	ctx, ctxCancel := context.WithCancel(c)
	_self := &ObserveHandler{
		id:        id,
		Ticker:    time.NewTicker(time.Millisecond * 60),
		dataChan:  make(chan rune, 64),
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
	monitor.AddTerm(id)
	go _self.notify()

	return _self
}

func (_self *ObserveHandler) GetIndex() int {
	return DefaultCmdHandler.GetIndex() + 2
}

func (_self *ObserveHandler) GetId() uint64 {
	return _self.id
}

func (_self *ObserveHandler) PassToClient(r rune) {
	_self.dataChan <- r

}
func (_self *ObserveHandler) notify() {
	var buf bytes.Buffer
	var err error
	for {
		select {
		case <-_self.C:
			dataStr := buf.String()
			if len(dataStr) <= 0 {
				continue
			}
			async.CommonWorkPool.Submit(func() {
				monitor.Subject.NotifyUpdateObservers(_self.GetId(), dataStr)
			})
			buf.Reset()

		case <-_self.ctx.Done():
			return
		case data := <-_self.dataChan:
			if data == utf8.RuneError {
				continue
			}

			if err = _self.runeToByte(data, func(b []byte) (err error) {
				_, err = buf.Write(b)
				return
			}); err != nil {
				return
			}

		}
	}

}

func (_self *ObserveHandler) PassToServer(r rune) bool {
	return true
}

func (_self *ObserveHandler) Close() {
	_self.ctxCancel()
	monitor.RemoveTerm(_self.GetId())
}

func (_self *ObserveHandler) runeToByte(r rune, write func([]byte) error) error {
	temp := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(temp, r)
	return write(temp)
}
