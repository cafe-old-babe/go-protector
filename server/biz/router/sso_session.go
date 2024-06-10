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
		routerGroup := group.Group("sso-session")
		{
			routerGroup.POST("/create/:authId", _ssoSession.CreateSession)
			routerGroup.GET("/connect/:ssoSessionId", _ssoSession.ConnectSession)
		}
	})

}

var _ssoSession ssoSession

type ssoSession struct {
	base.Router
}

func (_self ssoSession) CreateSession(c *gin.Context) {
	var authIdStr string
	authIdStr = c.Param("authId")
	authId, err := strconv.ParseUint(authIdStr, 10, 64)
	if err != nil {
		c_result.Err(c, err)
		return
	}
	var ssoSessionService service.SsoSession
	_self.MakeService(c, &ssoSessionService)
	result := ssoSessionService.CreateSession(authId)
	c_result.Result(c, result)
}

func (_self ssoSession) ConnectSession(c *gin.Context) {
	var req dto.ConnectSessionReq
	if err := c.ShouldBind(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var ssoSessionService service.SsoSession
	_self.MakeService(c, &ssoSessionService)
	if err := ssoSessionService.ConnectSession(&req); err != nil {
		c_result.Err(c, err)
		return
	}

}
