package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/database/condition"
	"go-protector/server/internal/utils/async"
	"gorm.io/gorm"
	"strconv"
)

type ApproveRecord struct {
	base.Service
}

// Page 分页查询
func (_self *ApproveRecord) Page(req *dto.ApproveRecordPageReq) (res *base.Result) {
	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var slice []entity.ApproveRecord
	tx := _self.GetDB().Scopes(
		condition.Paginate(req),
		condition.Like("ApproveUser.username", req.ApproveUsername),
		condition.Like("ApplicantUser.username", req.ApplicantUsername),
	)
	count, err := _self.Count(tx.Joins("ApproveUser").Joins("ApplicantUser").Find(&slice))
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	res = base.ResultPage(slice, req, count)

	return
}

// Insert 新增 res.Data *entity.ApproveRecord
func (_self *ApproveRecord) Insert(insertDTO *dto.ApproveRecordInsertDTO) (res *base.Result) {
	if insertDTO == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	if err := binding.Validator.ValidateStruct(insertDTO); err != nil {
		res = base.ResultFailureErr(err)
	}
	record := entity.ApproveRecord{
		SessionId:        insertDTO.SessionId,
		ApplicantId:      insertDTO.ApplicantId,
		ApplicantContent: insertDTO.ApplicantContext,
		ApproveUserId:    insertDTO.ApproveUserId,
		ApproveStatus:    consts.ApproveStatusUnprocessed,
		ApproveType:      insertDTO.ApproveType,
		ApproveBindId:    insertDTO.Id,
		Timeout:          insertDTO.Timeout,
	}
	err := _self.GetDB().Create(&record).Error
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	if insertDTO.Timeout > 0 {
		_, _ = async.NewDelayTask(strconv.Itoa(int(record.ID)), insertDTO.Timeout, _self.GetContext(), func(ctx context.Context) {
			//var recordService ApproveRecord
			//recordService.Make(ctx)

		})
	}
	res = base.ResultSuccess(&record)
	return
}

// DoApprove 审批
func (_self *ApproveRecord) DoApprove(doApproveDTO *dto.DoApproveDTO) (res *base.Result) {
	if doApproveDTO == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}

	if err := binding.Validator.ValidateStruct(doApproveDTO); err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	if doApproveDTO.ApproveUserId <= 0 {
		doApproveDTO.ApproveUserId = current.GetUserId(_self.GetContext())
	}
	if doApproveDTO.ApproveUserId <= 0 {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var record entity.ApproveRecord
	if err := _self.GetDB().First(&record, "id = ? and approve_status = ?", doApproveDTO.Id, consts.ApproveStatusUnprocessed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("记录不存或已被处理")
		}

		res = base.ResultFailureErr(err)
		return
	}

	// 校验状态
	switch doApproveDTO.ApproveStatus {
	case consts.ApproveStatusPass, consts.ApproveStatusReject:
		if record.ApproveUserId != doApproveDTO.ApproveUserId {
			res = base.ResultFailureErr(c_error.ErrAuthFailure)
			return
		}
	case consts.ApproveStatusTimeout, consts.ApproveStatusCancel:
		if record.ApplicantId != doApproveDTO.ApproveUserId {
			res = base.ResultFailureErr(c_error.ErrAuthFailure)
			return
		}
	default:
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	var update entity.ApproveRecord
	update.ID = doApproveDTO.Id
	update.ApproveStatus = doApproveDTO.ApproveStatus
	update.ApproveContent = doApproveDTO.ApproveContent

	tx := _self.GetDB().Updates(update)
	if err := tx.Error; err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	if tx.RowsAffected != 1 {
		res = base.ResultFailureErr(errors.New("处理失败,请稍后重试"))
		return

	}
	res = base.ResultSuccessMsg()

	return
}
