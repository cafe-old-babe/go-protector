package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/biz/service/iface"
	"go-protector/server/internal/approve"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_logger"
	"go-protector/server/internal/database/condition"
	"go-protector/server/internal/ssh/notify"
	"go-protector/server/internal/utils/async"
	"gorm.io/gorm"
	"strconv"
)

func init() {
	_ = approve.RegisterTypeInfo(consts.ApproveTypSsoOperation, "单点登录操作", func(ctx context.Context) error {
		record := ctx.Value("record").(entity.ApproveRecord)
		log := c_logger.GetLoggerByCtx(ctx)
		marshal, _ := json.Marshal(record)
		log.Debug("callback record: %s", string(marshal))
		notify.ApproveManager.Notify(record)
		return nil
	})
	iface.RegisterApproveRecordService(&ApproveRecord{})
}

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
		condition.Like(table_name.ApproveRecord+".work_num", req.WorkNum),
		condition.Like("ApproveUser.username", req.ApproveUsername),
		condition.Like("ApplicantUser.username", req.ApplicantUsername),
		func(db *gorm.DB) *gorm.DB {
			if user, ok := current.GetUser(_self.GetContext()); ok {
				if user.IsAdmin {
					return db
				}
				db = db.Where("( "+table_name.ApproveRecord+".applicant_id = @userId or "+
					table_name.ApproveRecord+".applicant_id = @userId )", sql.Named("userId", user.ID))
			}
			return db
		},
	)
	count, err := _self.Count(tx.Joins("ApproveUser").Joins("ApplicantUser").Order(table_name.ApproveRecord + ".created_at desc").Find(&slice))
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
		return
	}
	sysUser, err := dao.SysUser.FindUserByDTO(_self.GetDB(), dto.FindUserDTO{
		ID: insertDTO.ApproveUserId,
	})
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	record := entity.ApproveRecord{
		SessionId:        insertDTO.SessionId,
		ApplicantId:      insertDTO.ApplicantId,
		ApplicantContent: insertDTO.ApplicantContent,
		ApproveUserId:    insertDTO.ApproveUserId,
		ApproveStatus:    consts.ApproveStatusUnprocessed,
		ApproveType:      insertDTO.ApproveType,
		ApproveBindId:    insertDTO.ApproveBindId,
		Timeout:          insertDTO.Timeout,
	}

	err = _self.GetDB().Create(&record).Error
	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	record.ApproveUser = sysUser
	if insertDTO.Timeout > 0 {
		_self.Set("timeout", &dto.DoApproveDTO{
			Id:             record.ID,
			ApproveStatus:  consts.ApproveStatusTimeout,
			ApproveContent: "超时",
			ApproveUserId:  insertDTO.ApplicantId,
		})

		_, _ = async.NewDelayTask(strconv.Itoa(int(record.ID)), insertDTO.Timeout, _self.GetContext(), func(ctx context.Context) {
			if timeout, ok := ctx.Value("timeout").(*dto.DoApproveDTO); ok {
				var recordService ApproveRecord
				recordService.Make(ctx)
				recordService.DoApprove(timeout)
			}

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
	case consts.ApproveStatusReject:
		if len(doApproveDTO.ApproveContent) <= 0 {
			res = base.ResultFailureErr(errors.New("请填写拒绝通过原因"))
			return
		}
		fallthrough
	case consts.ApproveStatusPass:
		if len(doApproveDTO.ApproveContent) <= 0 {
			doApproveDTO.ApproveContent = "通过(系统默认生成)"
		}
		if record.ApproveUserId != doApproveDTO.ApproveUserId {
			res = base.ResultFailureErr(c_error.ErrAuthFailure)
			return
		}
	case consts.ApproveStatusCancel:
		if len(doApproveDTO.ApproveContent) <= 0 {
			doApproveDTO.ApproveContent = "用户撤回审批(系统默认生成)"
		}
		fallthrough
	case consts.ApproveStatusTimeout:
		if len(doApproveDTO.ApproveContent) <= 0 {
			doApproveDTO.ApproveContent = "超时未审批(系统默认生成)"
		}
		if record.ApplicantId != doApproveDTO.ApproveUserId {
			res = base.ResultFailureErr(c_error.ErrAuthFailure)
			return
		}
	default:
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
	}
	if record.Timeout > 0 {
		if err := async.CancelDelayTask(strconv.FormatUint(record.ID, 10)); err != nil {
			_self.GetLogger().Warn("CancelDelayTask err: %s")
		}
	}
	var update entity.ApproveRecord
	update.ID = doApproveDTO.Id
	update.ApproveStatus = doApproveDTO.ApproveStatus
	update.ApproveContent = doApproveDTO.ApproveContent

	tx := _self.GetDB().Updates(&update)
	if err := tx.Error; err != nil {
		res = base.ResultFailureErr(err)
		return
	}
	if tx.RowsAffected != 1 {
		res = base.ResultFailureErr(errors.New("处理失败,请稍后重试"))
		return
	}
	// 回调
	switch doApproveDTO.ApproveStatus {
	case consts.ApproveStatusPass, consts.ApproveStatusReject, consts.ApproveStatusTimeout:
		record.ApproveStatus = update.ApproveStatus
		record.ApproveStatusText = record.ApproveStatus.String()
		record.ApproveContent = update.ApproveContent
		ctx := context.WithValue(context.Background(), "record", record)
		_ = approve.HandleCallback(record.ApproveType, ctx)
	default:

	}
	res = base.ResultSuccessMsg()

	return
}
