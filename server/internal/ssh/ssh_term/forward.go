package ssh_term

import (
	"bytes"
	"context"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/ssh/cmd"
	"slices"
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
	chain      []cmd.Handler
}

func NewTermForward(ws *base.WsContext, term *Terminal, chain ...cmd.Handler) *TermForward {
	ctx, ctxCancel := context.WithCancel(ws.GetContext())
	if chain == nil {
		chain = make([]cmd.Handler, 0)
	}
	slices.SortFunc(chain, func(a, b cmd.Handler) int {
		return a.GetIndex() - b.GetIndex()
	})
	forward := &TermForward{
		ws:         ws,
		term:       term,
		dataChan:   make(chan rune),
		ctx:        ctx,
		ctxCancel:  ctxCancel,
		timeTicker: time.NewTicker(time.Millisecond * 60),
		chain:      chain,
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
				var next bool
				for i := range _self.chain {
					if next = _self.chain[i].PassToClient(data); !next {
						break
					}
				}
				if next {
					_self.dataChan <- data
				}
			}

		}

	}
}

func (_self *TermForward) readChanToWriteWs() {
	var buf bytes.Buffer
	var err error
	for {
		select {
		case <-_self.timeTicker.C:
			dataStr := buf.String()
			if len(dataStr) <= 0 {
				continue
			}
			if err = _self.ws.Write(base.NewWsMsg(consts.MsgData, dataStr)); err != nil {
				return
			}
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
			for _, r := range []rune(wsMsg.Body) {
				var next bool
				for i := range _self.chain {
					if next = _self.chain[i].PassToServer(r); !next {
						break
					}
				}
				if !next {
					continue
				}
				if err = _self.runeToByte(r, func(b []byte) (err error) {
					_, err = _self.term.Write(b)
					return
				}); err != nil {
					return
				}
			}

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

func (_self *TermForward) runeToByte(r rune, write func([]byte) error) error {
	temp := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(temp, r)
	return write(temp)
}
