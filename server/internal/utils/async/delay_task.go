package async

import (
	"context"
	"errors"
	"go-protector/server/internal/custom/c_error"
	"sync"
	"time"
)

var delayTaskMap sync.Map

type DelayTask struct {
	id        string
	f         func(ctx context.Context)
	cancel    context.CancelFunc
	delayTime time.Duration
}

func (_self *DelayTask) start(c context.Context) {
	timer := time.NewTimer(_self.delayTime)
	select {
	case <-timer.C:
		_, loaded := delayTaskMap.LoadAndDelete(_self.id)
		if loaded {
			_self.f(c)
		}

	case <-c.Done():
		_, _ = delayTaskMap.LoadAndDelete(_self.id)
	}
}

// Cancel 取消执行
func (_self *DelayTask) Cancel() (err error) {
	value, loaded := delayTaskMap.LoadAndDelete(_self.id)
	if !loaded {
		err = errors.New("任务不存在或已被执行")
		return
	}
	delayTask := value.(*DelayTask)
	delayTask.cancel()
	return
}

// NewDelayTask 创建延迟任务
func NewDelayTask(id string, delaySecond int64, c context.Context, f func(context.Context)) (*DelayTask, error) {
	if len(id) <= 0 || f == nil || delaySecond <= 0 || c == nil {
		return nil, c_error.ErrParamInvalid
	}

	_, ok := delayTaskMap.Load(id)
	if ok {
		return nil, errors.New("任务重复")
	}

	ctx, cancelFunc := context.WithCancel(context.WithoutCancel(c))
	task := DelayTask{
		id:        id,
		f:         f,
		cancel:    cancelFunc,
		delayTime: time.Duration(delaySecond) * time.Second,
	}
	delayTaskMap.Store(id, &task)
	go task.start(ctx)
	return &task, nil
}
