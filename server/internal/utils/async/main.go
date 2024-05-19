package async

import (
	"go-protector/server/internal/custom/c_logger"
	"go.uber.org/zap"
)

var MainWork *Main

type Main struct {
	queue chan func()
}

func NewMain() *Main {
	_self := &Main{queue: make(chan func(), 2048)}
	go _self.run()

	return _self

}

func (_self *Main) AsyncRun(f func()) {
	if f == nil {
		return
	}
	_self.queue <- f

}

func (_self *Main) run() {
	for {
		f, ok := <-_self.queue
		if !ok {
			break
		}
		if f == nil {
			continue
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					c_logger.ErrorZap("main 发生异常", zap.Any("err", err))
				}
			}()
			f()
		}()

	}
}
