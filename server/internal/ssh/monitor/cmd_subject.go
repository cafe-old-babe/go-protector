package monitor

import (
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/ssh/monitor/observe"
	"sync"
)

var Subject *cmdSubject

func init() {
	Subject = &cmdSubject{
		RWMutex: sync.RWMutex{},
		termMap: make(map[uint64][]observe.Observer),
	}
}

type cmdSubject struct {
	sync.RWMutex
	termMap map[uint64][]observe.Observer
}

func (_self *cmdSubject) RegisterObserver(observer observe.Observer) error {
	_self.Lock()
	defer _self.Unlock()
	id := observer.GetSsoId()
	observers, ok := _self.termMap[id]
	if !ok {
		return c_error.ErrIllegalAccess
	}
	_self.termMap[id] = append(observers, observer)
	return nil
}

func (_self *cmdSubject) RemoveObserver(observer observe.Observer) {
	// 删除 observers
	_self.Lock()
	defer _self.Unlock()
	id := observer.GetSsoId()
	observers, ok := _self.termMap[id]
	if !ok {
		return
	}
	for i, v := range observers {
		if v.GetObId() == observer.GetObId() {
			observers = append(observers[:i], observers[i+1:]...)
			break
		}
	}
}

func (_self *cmdSubject) NotifyUpdateObservers(id uint64, str string) {
	_self.RLock()
	_self.RUnlock()
	observers, ok := _self.termMap[id]
	if !ok {
		return
	}
	for _, observer := range observers {
		observer.Update(str)
	}
}
func (_self *cmdSubject) NotifyCloseObservers(id uint64) {
	_self.RLock()
	_self.RUnlock()
	observers, ok := _self.termMap[id]
	if !ok {
		return
	}
	for _, observer := range observers {
		observer.Close()
	}
}

func (_self *cmdSubject) addTerm(id uint64) {
	_self.Lock()
	defer _self.Unlock()
	if _, ok := _self.termMap[id]; ok {
		return
	}
	_self.termMap[id] = make([]observe.Observer, 0)
	return
}

func (_self *cmdSubject) removeTerm(id uint64) {
	_self.NotifyCloseObservers(id)
	_self.Lock()
	_self.Unlock()
	delete(_self.termMap, id)
}

func AddTerm(id uint64) {
	Subject.addTerm(id)
}
func RemoveTerm(id uint64) {
	Subject.removeTerm(id)
}
