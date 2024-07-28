package async

import (
	"context"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/custom/c_logger"
	"go.uber.org/zap"
	"sync"
)

// 6-18	【实战】采集/拨测资源从账号-使用协程异步采集-掌握利用通道实现CSP模型、协程间通信
var CommonWork *Work

func init() {
	//c_logger.SetLogger(zap.New(zapcore.NewNopCore()))
}

type Work struct {
	queue     chan func()
	c         context.Context
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

	workName := _self.c.Value(consts.CtxKeyTraceId).(string)
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
					c_logger.ErrorZap("发生异常", zap.String("workName", workName), zap.Any("err", err))
				}
			}()
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
