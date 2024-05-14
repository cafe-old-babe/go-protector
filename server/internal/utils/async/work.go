package async

import (
	"context"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/custom/c_logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

func init() {
	c_logger.SetLogger(zap.New(zapcore.NewNopCore()))
}

type Work struct {
	queue     chan func()
	c         context.Context
	log       *c_logger.SelfLogger
	wait      sync.WaitGroup
	closeOnce sync.Once
}

func NewWork(name string, limit int) *Work {
	if limit <= 0 {
		limit = 2048
	}
	ctx := context.WithValue(context.Background(), consts.CtxKeyTraceId, name)
	_self := &Work{
		queue: make(chan func(), limit),
		c:     ctx,
		log:   c_logger.GetLoggerByCtx(ctx),
	}
	_self.wait.Add(1)
	go _self.start()
	return _self
}

func (_self *Work) Submit(f func()) {
	if f == nil {
		return
	}
	_self.queue <- f
}

func (_self *Work) start() {
	if _self.queue == nil {
		return
	}

	for {
		f, ok := <-_self.queue
		if !ok {
			_self.wait.Done()
			break
		}
		func() {
			if f == nil {
				return
			}
			defer func() {
				if err := recover(); err != nil {
					_self.log.Error("发生异常, %+v", err)
				}
				_self.log.Debug("执行结束")
			}()
			_self.log.Debug("执行开始")

			f()
		}()

	}
}

func (_self *Work) Close() {
	_self.closeOnce.Do(func() {
		close(_self.queue)
	})
}
func (_self *Work) Wait() {
	defer _self.wait.Wait()
	_self.Close()
}
