package cmd

import "sync"

type cursor struct {
	x, y  int
	mutex sync.Mutex
}

func (_self *cursor) MoveX(i int) int {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()

	_self.x += i
	if _self.x < 0 {
		_self.x = 0
	}
	return _self.x
}

func (_self *cursor) MoveY(i int) int {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.y += i
	if _self.y < 0 {
		_self.y = 0
	}
	return _self.y
}

func (_self *cursor) GetX() int {
	return _self.x
}

func (_self *cursor) GetY() int {
	return _self.y
}

func (_self *cursor) SetX(i int) {
	_self.x = i
}

func (_self *cursor) SetY(i int) {
	_self.y = i
}

func (_self *cursor) ResetX() {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.x = 0
}

func (_self *cursor) ResetY() {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.y = 0
}

func (_self *cursor) ResetCursor() {
	_self.mutex.Lock()
	defer _self.mutex.Unlock()
	_self.y = 0
	_self.x = 0

}
