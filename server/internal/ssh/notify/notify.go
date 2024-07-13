package notify

import (
	"errors"
	"go-protector/server/internal/base"
	"go-protector/server/internal/ws"
	"sync"
)

var notifyMap = map[uint64]ws.IWsWriter{}

var lock sync.RWMutex

func RegisterWriter(id uint64, writer ws.IWsWriter) {
	lock.Lock()
	defer lock.Unlock()
	notifyMap[id] = writer
}

func UnRegisterWriter(id uint64) {
	lock.Lock()
	defer lock.Unlock()
	delete(notifyMap, id)
}

func WriterById(id uint64, msg *base.WsMsg) (err error) {

	lock.RLock()
	defer lock.RUnlock()

	writer, ok := notifyMap[id]
	if !ok {
		err = errors.New("告警会话不存在或已经结束")
		return
	}
	err = writer.Write(msg)

	return
}
