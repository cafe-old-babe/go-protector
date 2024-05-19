package base

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_logger"
	"go-protector/server/internal/database"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

type IService interface {
	Make(ctx *gin.Context)
	MakeService(service ...IService)
	GetDB() *gorm.DB
	WithGoroutineDB()
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
func (_self *Service) GetDB() *gorm.DB {
	return _self.DB
}

// WithGoroutineDB 开启协程 处理数据库连接 防止 context canceled
func (_self *Service) WithGoroutineDB() {
	_self.Context = _self.Context.Copy()
	_self.DB = _self.DB.WithContext(context.Background())
	_self.Context.Set(consts.CtxKeyDB, _self.DB)

}

// SimpleSave 通用保存 会保存所有的字段，即使字段是零值 简单
func (_self *Service) SimpleSave(model schema.Tabler,
	check ...func() error) *Result {
	if len(check) > 0 {
		for _, f := range check {
			if err := f(); err != nil {
				return ResultFailureErr(err)
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
		return ResultFailureErr(err)
	}
	return ResultSuccess(model, message)
}

// SimpleDelByIds 通用删除方法 根据ids删除数据
func (_self *Service) SimpleDelByIds(req *IdsReq,
	check ...func() error) *Result {

	if len(check) > 0 {
		for _, f := range check {
			if err := f(); err != nil {
				return ResultFailureErr(err)
			}
		}
	}
	if req == nil || len(req.GetIds()) <= 0 || req.Value == nil {
		return ResultFailureErr(c_error.ErrParamInvalid)
	}
	if len(check) > 0 {
		for _, f := range check {
			if err := f(); err != nil {
				return ResultFailureErr(err)
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
		return ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.Logger.Error("删除失败,无删除记录")
		return ResultFailureMsg("删除失败")
	}
	return ResultSuccessMsg("删除成功")

}

// Count 获取总数
func (_self *Service) Count(db *gorm.DB) (count int64, err error) {
	if err = db.Error; err != nil {
		return
	}
	err = db.Limit(-1).Offset(-1).Count(&count).Error
	return
}
