package ws

import (
	"go-protector/server/internal/base"
	"go-protector/server/internal/current"
	"strconv"
	"sync"
)

type manager struct {
	group              map[string]map[string]*wsClient
	msgChan            chan *sendWsCliMsg
	regChan, unRegChan chan *wsClient
	sync.RWMutex
}

type sendWsCliMsg struct {
	groupId, id string
	*base.WsMsg
}
type wsClient struct {
	ICtxWsWriteCloser
	registerFunc []func(service base.IService)
}

func (_self *manager) start() {
	for {
		select {
		case msg := <-_self.msgChan:
			go _self.send(msg)
		case cli := <-_self.regChan:
			_self.register(cli)
		case cli := <-_self.unRegChan:
			_self.unRegister(cli)
		}
	}

}

func (_self *manager) send(msg *sendWsCliMsg) {
	if msg == nil {
		return
	}
	if len(msg.groupId) <= 0 {
		return
	}
	_self.RLock()
	defer _self.RUnlock()
	group, ok := _self.group[msg.groupId]
	if !ok {
		return
	}
	if len(group) <= 0 {
		return
	}
	if len(msg.id) <= 0 {
		// 组播
		for _, cli := range group {
			_ = cli.Write(msg.WsMsg)
		}
	} else {
		// 单播
		_ = group[msg.id].Write(msg.WsMsg)
	}

}

// register 注册
func (_self *manager) register(cli *wsClient) {
	user, ok := current.GetUser(cli.GetContext())
	if !ok || user.ID <= 0 || len(user.SessionId) <= 0 {
		return
	}
	id := user.SessionId
	groupId := strconv.FormatUint(user.ID, 10)
	if len(id) <= 0 || len(groupId) <= 0 {
		return
	}
	_self.Lock()
	defer _self.Unlock()
	if _, ok := _self.group[groupId]; !ok {
		_self.group[groupId] = make(map[string]*wsClient)
	}
	if oldCli, ok := _self.group[groupId][id]; ok {
		_ = oldCli.Close()
		delete(_self.group[groupId], id)
	}
	_self.group[groupId][id] = cli
	for _, f := range cli.registerFunc {
		f(cli)
	}

}

// unRegister 取消注册
func (_self *manager) unRegister(cli ICtxWsWriteCloser) {
	user, ok := current.GetUser(cli.GetContext())
	if !ok || user.ID <= 0 || len(user.SessionId) <= 0 {
		return
	}
	id := user.SessionId
	groupId := strconv.FormatUint(user.ID, 10)
	if len(id) <= 0 || len(groupId) <= 0 {
		return
	}
	_self.Lock()
	defer _self.Unlock()
	if _, ok := _self.group[groupId]; !ok {
		return
	}
	delete(_self.group[groupId], id)
}

var _manager manager
var once sync.Once

func init() {
	once.Do(func() {
		_manager = manager{
			group:     make(map[string]map[string]*wsClient),
			msgChan:   make(chan *sendWsCliMsg, 2048),
			regChan:   make(chan *wsClient, 2048),
			unRegChan: make(chan *wsClient, 2048),
		}
		go _manager.start()
	})

}

// RegisterWsCli 注册
func RegisterWsCli(ws ICtxWsWriteCloser, f ...func(service base.IService)) {
	if ws == nil {
		return
	}
	_manager.regChan <- &wsClient{
		ICtxWsWriteCloser: ws,
		registerFunc:      f,
	}
}

// UnRegisterWsCli 取消注册
func UnRegisterWsCli(ws ICtxWsWriteCloser) {
	if ws == nil {
		return
	}
	_manager.unRegChan <- &wsClient{
		ICtxWsWriteCloser: ws,
	}

}

// SendMsgByGroupId 组播
func SendMsgByGroupId(msg *base.WsMsg, groupId string) {
	if len(groupId) <= 0 || msg == nil {
		return
	}
	_manager.msgChan <- &sendWsCliMsg{
		groupId: groupId,
		WsMsg:   msg,
	}
}

// SendMsgById 单播
func SendMsgById(msg *base.WsMsg, groupId, id string) {
	if len(groupId) <= 0 || len(id) <= 0 || msg == nil {
		return
	}
	_manager.msgChan <- &sendWsCliMsg{
		groupId: groupId,
		WsMsg:   msg,
	}
}
