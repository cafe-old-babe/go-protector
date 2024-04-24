package base

import "github.com/gin-gonic/gin"

type Api struct {
}

func (_self Api) MakeService(c *gin.Context, service ...IService) {
	if len(service) <= 0 {
		return
	}
	for i := range service {
		service[i].Make(c)
	}
}
