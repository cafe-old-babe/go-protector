package base

import "github.com/gin-gonic/gin"

type Router struct {
}

func (_self Router) MakeService(c *gin.Context, service ...IService) {
	if len(service) <= 0 {
		return
	}
	for i := range service {
		service[i].Make(c)
	}
}
