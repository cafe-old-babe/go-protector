package approve

import (
	"context"
	"errors"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/utils/async"
	"sync"
)

var infoType = map[c_type.ApproveType]*TypeInfo{}
var lock sync.RWMutex

type TypeInfo struct {
	approveType c_type.ApproveType
	// ApproveTypeText 审批类型
	ApproveTypeText string
	// HandleCallback  审批结果在 ctx中 key: record; value: *entity.ApproveRecord
	HandleCallback func(ctx context.Context) error
}

func RegisterTypeInfo(approveType c_type.ApproveType, ApproveTypeText string, f func(ctx context.Context) error) (err error) {
	if len(approveType) <= 0 || len(ApproveTypeText) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	lock.Lock()
	defer lock.Unlock()
	infoType[approveType] = &TypeInfo{
		approveType:     approveType,
		ApproveTypeText: ApproveTypeText,
		HandleCallback:  f,
	}
	return
}

func GetTypeInfo(approveType c_type.ApproveType) *TypeInfo {

	lock.RLock()
	defer lock.RUnlock()
	return infoType[approveType]
}

func HandleCallback(approveType c_type.ApproveType, ctx context.Context) error {
	info := GetTypeInfo(approveType)
	if info == nil {
		return c_error.ErrIllegalAccess
	}

	if info.HandleCallback == nil {
		return errors.New("未注册回调方法")
	}

	async.CommonWorkPool.Submit(func() {
		_ = info.HandleCallback(context.WithoutCancel(ctx))
	})
	return nil

}
