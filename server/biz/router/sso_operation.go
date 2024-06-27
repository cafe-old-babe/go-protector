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

		routerGroup := group.Group("sso-operation")
		{
			routerGroup.POST("/page", _ssoOperation.Page)
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
