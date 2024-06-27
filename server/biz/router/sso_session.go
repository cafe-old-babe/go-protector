package router

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/service"
	"go-protector/server/internal/base"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_result"
	"strconv"
)

func init() {

	initRouterFunc = append(initRouterFunc, func(group *gin.RouterGroup) {
		group.GET("/ws/sso-session/connect/:ssoSessionId", _ssoSession.ConnectBySession)

		routerGroup := group.Group("sso-session")
		{
			routerGroup.POST("/create/:authId", _ssoSession.CreateSession)
			routerGroup.POST("/page", _ssoSession.Page)
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

func (_self ssoSession) ConnectBySession(c *gin.Context) {
	var req dto.ConnectBySessionReq
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	if str := c.Param("ssoSessionId"); len(str) <= 0 {
		c_result.Err(c, c_error.ErrParamInvalid)
		return
	} else {
		if req.Id, err = strconv.ParseUint(str, 10, 64); err != nil {
			c_result.Err(c, c_error.ErrParamInvalid)
			return
		}
	}

	var ssoSessionService service.SsoSession
	_self.MakeService(c, &ssoSessionService)
	if err = ssoSessionService.ConnectBySession(&req); err != nil {
		c_result.Err(c, err)
		return
	}

}

func (_self ssoSession) Page(c *gin.Context) {
	var req dto.SsoSessionPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c_result.Err(c, err)
		return
	}
	var sessionService service.SsoSession
	_self.MakeService(c, &sessionService)

	c_result.Result(c, sessionService.Page(&req))
}
