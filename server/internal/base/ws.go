package base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"net/http"
	"strconv"
	"strings"
	"syscall"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  syscall.Getpagesize(),
	WriteBufferSize: syscall.Getpagesize(),
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsContext struct {
	conn *websocket.Conn
	user *current.User
	Service
}

func (_self *WsContext) Value(key string) (value any) {
	value = _self.GetContext().Value(key)
	return
}

// Write 写入消息
func (_self *WsContext) Write(msg *WsMsg) error {
	return _self.conn.WriteMessage(websocket.TextMessage, msg.ToBytes())
}

// ReadMsg 读取消息
func (_self *WsContext) ReadMsg() (wsMsg *WsMsg, err error) {
	var bytes []byte
	if _, bytes, err = _self.conn.ReadMessage(); err != nil {
		return
	}
	if len(bytes) <= 0 {
		return
	}
	if wsMsg, err = ParseWsMsg(string(bytes)); err != nil {
		_self.GetLogger().Error("parse ws msg err: %v", err)
		err = nil
		return
	}

	return
}

func (_self *WsContext) Close() error {
	return _self.conn.Close()
}

func (_self *WsContext) GetUser() *current.User {
	return _self.user
}

func Upgrade(service IService) (*WsContext, error) {
	var user *current.User
	var ok bool
	var err error
	if user, ok = current.GetUser(service.GetContext()); !ok {
		return nil, c_error.ErrParamInvalid
	}
	var ginCtx *gin.Context
	if ginCtx, err = service.GetGinCtx(); err != nil {
		return nil, err
	}
	var ws *websocket.Conn
	ws, err = Upgrader.Upgrade(ginCtx.Writer, ginCtx.Request, nil)
	if err != nil {
		return nil, err
	}
	wsContext := &WsContext{
		conn: ws,
		user: user,
	}
	service.MakeService(wsContext)
	return wsContext, nil
}

type WsMsg struct {
	MsgNum c_type.MsgType
	Body   string
}

var msgFmt = "%d %s"

func (_self *WsMsg) ToBytes() []byte {
	sprintf := fmt.Sprintf(msgFmt, _self.MsgNum, _self.Body)
	return []byte(sprintf)
}

// ParseWsMsg 转换
func ParseWsMsg(str string) (msg *WsMsg, err error) {
	if len(str) <= 0 {
		return
	}
	i := strings.Index(str, " ")
	if i <= 0 {
		return
	}
	numStr := str[:i]
	var atoi int
	if atoi, err = strconv.Atoi(numStr); err != nil {
		return
	}
	msg = new(WsMsg)
	msg.MsgNum = c_type.MsgType(atoi)
	msg.Body = str[i+1:]

	return
}

func NewWsMsg(t c_type.MsgType, body string) *WsMsg {
	return &WsMsg{
		MsgNum: t,
		Body:   body,
	}
}
