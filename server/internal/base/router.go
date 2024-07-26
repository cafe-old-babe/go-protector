package base

import (
	"context"
)

type Router struct {
}

// 4-8	【实战】字典类型管理接口开发之发现问题解决问题-掌握架构抽象能力
func (_self Router) MakeService(c context.Context, service ...IService) {
	if len(service) <= 0 {
		return
	}
	for i := range service {
		service[i].Make(c)
	}
}
