package async

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

// 6-19	【实战】采集/拨测资源从账号-并行异步采集-掌握如何提高CSP模型吞吐量；掌握闭包基础概念
var CommonWorkPool *WorkPool

var num atomic.Int32

const workNameFmt = "%s-%d-%d"

type WorkPool struct {
	workSlice []*Work
	workSize  int32
	cur       atomic.Int32
}

// NewWorkPool create
// workSize work个数
// 每个work的workLimit
func NewWorkPool(name string, workSize, workLimit int) *WorkPool {

	if workSize <= 0 {
		// 获取系统协程数
		cup := runtime.NumCPU()
		workSize = cup << 2
	}
	if workLimit <= 0 {
		workLimit = 2048
	}
	poolNum := num.Add(1)

	if len(name) <= 0 {
		name = "work-pool"
	}
	_self := &WorkPool{
		workSlice: make([]*Work, workSize),
		workSize:  int32(workSize),
	}
	for i := 0; i < workSize; i++ {
		_self.workSlice[i] = NewWork(fmt.Sprintf(workNameFmt, name, poolNum, i), workLimit)
	}

	return _self
}

func (_self *WorkPool) Submit(f func()) {

	index := _self.cur.Add(1)
	if cur := _self.cur.Load(); cur >= _self.workSize {
		_self.cur.CompareAndSwap(cur, 0)
	}
	work := _self.workSlice[index-1]
	work.Submit(f)
}

func (_self *WorkPool) Close() {

	for i := range _self.workSlice {
		_self.workSlice[i].Close()
	}
}

func (_self *WorkPool) Wait() {
	_self.Close()
	for i := range _self.workSlice {
		_self.workSlice[i].Wait()
	}

}
