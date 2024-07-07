package c_structure

import "sync"

type SafeStack[T any] struct {
	stack []T
	sync.Mutex
}

func (_self *SafeStack[T]) Push(data T) {
	_self.Lock()
	defer _self.Unlock()
	_self.stack = append(_self.stack, data)
}

func (_self *SafeStack[T]) Pop() (data T, exists bool) {
	if exists = !_self.IsEmpty(); !exists {
		return
	}
	_self.Lock()
	defer _self.Unlock()

	i := len(_self.stack) - 1
	data = _self.stack[i]
	_self.stack = _self.stack[:i]
	return

}

func (_self *SafeStack[T]) Top() (data T, exists bool) {

	if exists = !_self.IsEmpty(); !exists {
		return
	}
	_self.Lock()
	defer _self.Unlock()
	data = _self.stack[len(_self.stack)-1]
	return
}

func (_self *SafeStack[T]) IsEmpty() bool {
	_self.Lock()
	defer _self.Unlock()
	return len(_self.stack) <= 0
}

func (_self *SafeStack[T]) Clear() {
	_self.Lock()
	defer _self.Unlock()
	_self.stack = _self.stack[:0]
	return
}
