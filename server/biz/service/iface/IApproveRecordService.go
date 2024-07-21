package iface

import (
	"context"
	"go-protector/server/biz/model/dto"
	"go-protector/server/internal/base"
	"reflect"
)

type IApproveRecordService interface {
	Make(c context.Context)
	// Insert 新增审批
	Insert(insertDTO *dto.ApproveRecordInsertDTO) (res *base.Result)
	// DoApprove 处理审批
	DoApprove(doApproveDTO *dto.DoApproveDTO) (res *base.Result)
}

var approveRecordService IApproveRecordService

func RegisterApproveRecordService(impl IApproveRecordService) {
	approveRecordService = impl
}

func ApproveRecordService(c context.Context) IApproveRecordService {

	value := reflect.New(reflect.ValueOf(approveRecordService).Elem().Type())
	//value := reflect.New(reflect.Indirect(reflect.ValueOf(approveRecordService)).Type())
	recordService := value.Interface().(IApproveRecordService)
	recordService.Make(c)
	return recordService
}
