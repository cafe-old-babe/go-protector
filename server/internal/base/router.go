package base

import (
	"context"
)

type Router struct {
}

func (_self Router) MakeService(c context.Context, service ...IService) {
	if len(service) <= 0 {
		return
	}
	for i := range service {
		service[i].Make(c)
	}
}
