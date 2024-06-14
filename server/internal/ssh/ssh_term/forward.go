package ssh_term

import (
	"bytes"
	"context"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"time"
	"unicode/utf8"
)

type TermForward struct {
	ws         *base.WsContext
	term       *Terminal
	dataChan   chan rune
	ctx        context.Context
	ctxCancel  context.CancelFunc
	timeTicker *time.Ticker
}

func NewTermForward(ws *base.WsContext, term *Terminal) *TermForward {
	ctx, ctxCancel := context.WithCancel(ws.GetContext())

	forward := &TermForward{
		ws:         ws,
		term:       term,
		dataChan:   make(chan rune),
		ctx:        ctx,
		ctxCancel:  ctxCancel,
		timeTicker: time.NewTicker(time.Millisecond * 60),
	}

	return forward
}

func (_self *TermForward) Start() {
	go _self.readTermToWriteChan()
	go _self.readChanToWriteWs()

}

func (_self *TermForward) Stop() {
	if _self.ws != nil {
		_ = _self.ws.Close()
	}
	_self.timeTicker.Stop()
	_self.ctxCancel()
}

func (_self *TermForward) readTermToWriteChan() {
	var data rune
	var i int
	var err error
	for {
		select {
		case <-_self.ctx.Done():
			return
		default:
			if data, i, err = _self.term.ReadRune(); err != nil {
				return
			}
			if i > 0 {
				_self.dataChan <- data
			}

		}

	}
}

func (_self *TermForward) readChanToWriteWs() {
	var buf bytes.Buffer
	for {
		select {
		case <-_self.timeTicker.C:
			dataStr := buf.String()
			if len(dataStr) <= 0 {
				continue
			}
			_ = _self.ws.Write(base.NewWsMsg(consts.MsgData, dataStr))
			buf.Reset()

		case <-_self.ctx.Done():
			return
		case data := <-_self.dataChan:
			temp := make([]byte, utf8.RuneLen(data))
			utf8.EncodeRune(temp, data)
			buf.Write(temp)

		}
	}
}

// ReadWsToWriteTerm 读取ws 写入 term
func (_self *TermForward) ReadWsToWriteTerm() (err error) {
	if err = _self.ws.Write(base.NewWsMsg(consts.MsgConnected, "")); err != nil {
		return
	}
	var wsMsg *base.WsMsg
	for {

		if wsMsg, err = _self.ws.ReadMsg(); err != nil {
			return
		}
		if wsMsg == nil {
			continue
		}
		switch wsMsg.MsgNum {
		case consts.MsgData:
			_, err = _self.term.Write([]byte(wsMsg.Body))
		case consts.MsgClose:
			return
		default:
			return
		}
		if err != nil {
			return
		}

	}
}
