package notify

import (
	"go-protector/server/biz/model/entity"
	"sync"
)

var ApproveManager *manager

type manager struct {
	approveMap map[uint64]func(record entity.ApproveRecord)
	mutex      sync.Mutex
}

// Subscribe 订阅
func (_self *manager) Subscribe(id uint64, f func(record entity.ApproveRecord)) {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.approveMap[id] = f
}

// UnSubscribe 取消订阅
func (_self *manager) UnSubscribe(id uint64) {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	delete(_self.approveMap, id)
}

// Notify 通知
func (_self *manager) Notify(record entity.ApproveRecord) {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	if f := _self.approveMap[record.ID]; f != nil {
		f(record)
		delete(_self.approveMap, record.ID)
	}
}

func init() {
	ApproveManager = &manager{
		approveMap: make(map[uint64]func(record entity.ApproveRecord)),
	}

}
