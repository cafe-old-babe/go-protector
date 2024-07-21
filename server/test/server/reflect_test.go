package server

import (
	"context"
	"fmt"
	"go-protector/server/biz/service"
	"go-protector/server/biz/service/iface"
	"testing"
)

func TestName(t *testing.T) {
	iface.RegisterApproveRecordService(&service.ApproveRecord{})
	recordService := iface.ApproveRecordService(context.Background())
	fmt.Printf("--%v\n", recordService)
}
