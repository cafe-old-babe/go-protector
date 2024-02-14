package base

import (
	"github.com/gin-gonic/gin"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_logger"
	"go-protector/server/core/database"
	"go-protector/server/models/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

type IService interface {
	Make(ctx *gin.Context)
	MakeService(service ...IService)
}

type Service struct {
	Logger  *c_logger.SelfLogger
	DB      *gorm.DB
	Context *gin.Context
}

func (_self *Service) Make(c *gin.Context) {
	_self.DB = database.GetDB(c)
	_self.Logger = c_logger.GetLogger(c)
	_self.Context = c
}

func (_self *Service) MakeService(service ...IService) {
	if len(service) <= 0 {
		return
	}
	for i := range service {
		service[i].Make(_self.Context)
	}
}

// CommonSave 通用保存 会保存所有的字段，即使字段是零值
func (_self *Service) CommonSave(model schema.Tabler,
	check ...func() error) *dto.Result {
	if len(check) > 0 {
		for _, f := range check {
			if err := f(); err != nil {
				return dto.ResultFailureErr(err)
			}
		}
	}
	// 反射获取 model的ID属性
	idValue := reflect.ValueOf(model).Elem().FieldByName("ID")
	var err error
	var message string
	if idValue.IsZero() {
		// 判断ID是否为空
		err = _self.DB.Create(model).Error
		message = "新增成功"
	} else {
		//更新属性，只会更新非零值的字段
		err = _self.DB.Updates(model).Error
		message = "更新成功"
	}
	if err != nil {
		return dto.ResultFailureErr(err)
	}
	return dto.ResultSuccess(model, message)
}

// CommonDelByIds 通用删除方法 根据ids删除数据
func (_self *Service) CommonDelByIds(req *dto.IdsReq,
	check ...func() error) *dto.Result {

	if len(check) > 0 {
		for _, f := range check {
			if err := f(); err != nil {
				return dto.ResultFailureErr(err)
			}
		}
	}
	if req == nil || len(req.GetIds()) <= 0 || req.Value == nil {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	if len(check) > 0 {
		for _, f := range check {
			if err := f(); err != nil {
				return dto.ResultFailureErr(err)
			}
		}
	}
	var tx *gorm.DB
	if req.Unscoped {
		tx = _self.DB.Unscoped()
	} else {
		tx = _self.DB
	}
	result := tx.Delete(req.Value, req.GetIds())
	if result.Error != nil {
		return dto.ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.Logger.Error("删除失败,无删除记录")
		return dto.ResultFailureMsg("删除失败")
	}
	return dto.ResultSuccessMsg("删除成功")

}
