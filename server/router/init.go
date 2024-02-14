package router

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	initSysUser(routerGroup)
	initSystem(routerGroup)
	initSysDict(routerGroup)
}
