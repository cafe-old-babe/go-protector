package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/commons/base"
)

var UserApi = new(User)

type User struct {
	base.Api
}

func (_self User) Login(c *gin.Context) {

}

func (_self User) Logout(c *gin.Context) {

}
