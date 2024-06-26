package c_structure

import "sync"

type SafeNonBlockingQueue[T any] struct {
	queue []T
	mutex sync.Mutex
}

func (_self *SafeNonBlockingQueue[T]) Push(data T) {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.queue = append(_self.queue, data)
}

func (_self *SafeNonBlockingQueue[T]) Pop() (data T, exists bool) {
	if exists = !_self.IsEmpty(); !exists {
		return
	}

	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	data = _self.queue[0]
	_self.queue = _self.queue[1:]
	return
}

func (_self *SafeNonBlockingQueue[T]) PeekHand() (data T, exists bool) {
	if exists = !_self.IsEmpty(); !exists {
		return
	}

	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	data = _self.queue[0]
	return
}

func (_self *SafeNonBlockingQueue[T]) PopAll(f func(T)) {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	var data T

	for len(_self.queue) > 0 {
		data = _self.queue[0]
		_self.queue = _self.queue[1:]
		f(data)
	}
	return
}

func (_self *SafeNonBlockingQueue[T]) IsEmpty() bool {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	return len(_self.queue) <= 0
}

func (_self *SafeNonBlockingQueue[T]) Clear() {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.queue = _self.queue[:0]
	return
}
