package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_result"
	"strconv"
)

func init() {

	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {

		routerGroup := group.Group("sso-operation")
		{
			routerGroup.POST("/page", _ssoOperation.Page)
			routerGroup.POST("/page/:ssoSessionId", _ssoOperation.PageBySsoId)
		}
	})

}

var _ssoOperation ssoOperation

type ssoOperation struct {
	base.Router
}

func (_self ssoOperation) Page(c *gin.Context) {
	var req dto.SsoOperationPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var operationService service.SsoOperation
	_self.MakeService(c, &operationService)

	c_result.Result(c, operationService.Page(&req))
}

func (_self ssoOperation) PageBySsoId(c *gin.Context) {

	params := c.Param("ssoSessionId")
	ssoSessionId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		c_result.Err(c, err)

		return
	}
	var req dto.SsoOperationPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	req.SsoSessionId = ssoSessionId

	var operateService service.SsoOperation
	_self.MakeService(c, &operateService)
	c_result.Result(c, operateService.PageBySsoId(&req))

}
