package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
)

func init() {
	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		routerGroup := group.Group("approve-cmd")
		{
			routerGroup.POST("/page", _approveCmd.Page)
			routerGroup.POST("/save", _approveCmd.Save)
			routerGroup.POST("/delete", _approveCmd.Delete)
		}
	})
}

var _approveCmd approveCmd

type approveCmd struct {
	base.Router
}

func (_self approveCmd) Page(c *gin.Context) {
	var req dto.ApproveCmdPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var approveService service.ApproveCmd
	_self.MakeService(c, &approveService)
	c_result.Result(c, approveService.Page(&req))

}

func (_self approveCmd) Save(c *gin.Context) {
	var req entity.ApproveCmd
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var approveService service.ApproveCmd
	_self.MakeService(c, &approveService)
	req.ID = 0
	res := approveService.SimpleSave(&req, func() error {
		return approveService.SaveCheck(req)
	})
	c_result.Result(c, res)
}

func (_self approveCmd) Delete(c *gin.Context) {
	var idsReq base.IdsReq
	err := c.ShouldBindJSON(&idsReq)
	if err != nil {
		c_result.Err(c, err)
		return
	}
	var approveService service.ApproveCmd
	_self.MakeService(c, &approveService)
	idsReq.Value = &entity.ApproveCmd{}
	result := approveService.SimpleDelByIds(&idsReq)
	c_result.Result(c, result)
}
