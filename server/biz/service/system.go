package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"go-protector/server/internal/ws"
	"strconv"
)

type System struct {
	base.Service
}

func MakeSystem(c *gin.Context) *System {
	var self System
	self.Make(c)
	return &self
}

func (_self *System) Bus() (err error) {
	wsContext, err := base.Upgrade(&_self.Service)
	if err != nil {
		return
	}
	defer func() {
		if wsContext != nil {
			ws.UnRegisterWsCli(wsContext)
			_ = wsContext.Close()
		}
	}()

	ws.RegisterWsCli(wsContext, func(service base.IService) {
		var approveService ApproveRecord
		approveService.Make(service.GetContext())
		var count int64
		var currentUser *current.User
		var ok bool
		if currentUser, ok = current.GetUser(service.GetContext()); !ok {
			approveService.GetLogger().Error("未找到当前用户信息")
			return
		}
		if count, err = approveService.FindCountUnprocessed(currentUser.ID); err != nil {
			approveService.GetLogger().Error("FindCountUnprocessed err: %v", err)
			return
		}
		if count <= 0 {
			return
		}

		groupId := strconv.FormatUint(currentUser.ID, 10)

		msg := base.NewWsMsg(consts.MsgApprove,
			fmt.Sprintf("%s 您好, 你有 %d 待处理的审批",
				currentUser.UserName, count))
		ws.SendMsgByGroupId(msg, groupId)
	})
	for {
		if _, err = wsContext.ReadMsg(); err != nil {
			return
		}
	}

}
