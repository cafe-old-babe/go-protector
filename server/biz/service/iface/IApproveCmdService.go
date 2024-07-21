package iface

import (
	"context"
	"reflect"
)

type IApproveCmdService interface {
	Make(c context.Context)
	GetApproveCmdSliceByAssetId(assetId uint64) (cmdSlice []string, err error)
}

var approveCmdService IApproveCmdService

func RegisterApproveCmdService(impl IApproveCmdService) {
	approveCmdService = impl
}

func ApproveCmdService(c context.Context) IApproveCmdService {
	value := reflect.New(reflect.Indirect(reflect.ValueOf(approveCmdService)).Type())
	recordService := value.Interface().(IApproveCmdService)
	recordService.Make(c)
	return recordService
}
