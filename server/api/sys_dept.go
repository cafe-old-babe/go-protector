package api

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/base"
)

var SysDeptApi sysDept

type sysDept struct {
	base.Api
}

// Page 分页查询
func (_self sysDept) Page(c *gin.Context) {

	return
}

// Delete 删除
func (_self sysDept) Delete(c *gin.Context) {

}

// Save 保存
func (_self sysDept) Save(c *gin.Context) {

}
