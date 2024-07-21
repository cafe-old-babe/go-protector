package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("approve-record")
		{
			routerGroup.POST("/page", _approveRecord.Page)
			routerGroup.POST("/doApprove", _approveRecord.doApprove)
		}
	})
}

var _approveRecord approveRecord

type approveRecord struct {
	base.Router
}

func (_self approveRecord) Page(c *gin.Context) {
	var req dto.ApproveRecordPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}

	var approveService service.ApproveRecord
	_self.MakeService(c, &approveService)
	c_result.Result(c, approveService.Page(&req))
}

func (_self approveRecord) doApprove(c *gin.Context) {
	var req dto.DoApproveDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var approveService service.ApproveRecord
	_self.MakeService(c, &approveService)
	c_result.Result(c, approveService.DoApprove(&req))
}
