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

// 4-8	【实战】字典类型管理接口开发之发现问题解决问题-掌握架构抽象能力
type IService interface {
	Make(ctx context.Context)
	MakeService(service ...IService)
	GetDB() *gorm.DB
	WithGoroutineDB()
	GetContext() context.Context
	GetLogger() *c_logger.SelfLogger
	GetGinCtx() (*gin.Context, error)
	Set(k, v any)
	Begin()
}

type Service struct {
	logger *c_logger.SelfLogger
	db     *gorm.DB
	ctx    context.Context
}

func (_self *Service) Make(c context.Context) {
	if ginCtx, ok := c.(*gin.Context); ok {
		c = context.WithValue(c, "ginCtx", ginCtx)
	}
	_self.ctx = c
	// 死锁
	// 7-6	【实战】解决通过钩子函数更新授权冗余数据造成死锁的问题-掌握排查死锁的思路
	// _self.db = database.GetDB(c)
	// 使用已有的db
	if db, ok := _self.ctx.Value(consts.CtxKeyDB).(*gorm.DB); ok {
		_self.db = db
	} else {
		_self.db = database.GetDB(c)
		_self.Set(consts.CtxKeyDB, _self.db)
	}

}

func (_self *Service) MakeService(service ...IService) {
	if len(service) <= 0 {
		return
	}
	//_self.ctx = context.WithValue(_self.GetContext(), consts.CtxKeyDB, _self.GetDB())
	for i := range service {
		if service[i] == nil {
			service[i] = new(Service)
		}
		service[i].Make(_self.GetContext())
	}
}

func (_self *Service) GetDB() *gorm.DB {
	return _self.db
}

// WithGoroutineDB 开启协程 处理数据库连接 防止 context canceled
func (_self *Service) WithGoroutineDB() {
	_self.ctx = context.WithoutCancel(_self.GetContext())
	_self.db = _self.db.WithContext(_self.GetContext())
	_self.ctx = context.WithValue(_self.GetContext(), consts.CtxKeyDB, _self.GetDB())

}

func (_self *Service) GetContext() context.Context {
	return _self.ctx
}

func (_self *Service) GetLogger() *c_logger.SelfLogger {
	if _self.logger == nil {
		_self.logger = c_logger.GetLoggerByCtx(_self.GetContext())
	}
	return _self.logger
}

func (_self *Service) GetGinCtx() (*gin.Context, error) {
	value := _self.GetContext().Value("ginCtx")
	if ginCtx, ok := value.(*gin.Context); ok {
		if ginCtx == nil {
			return nil, c_error.ErrIllegalAccess
		}
		return ginCtx, nil
	}

	return nil, c_error.ErrIllegalAccess

}

func (_self *Service) Set(k, v any) {
	_self.ctx = context.WithValue(_self.GetContext(), k, v)
}

func (_self *Service) Begin() {
	_self.db = _self.db.Begin()
	_self.Set(consts.CtxKeyDB, _self.db)
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
		err = _self.db.Create(model).Error
		message = "新增成功"
	} else {
		//更新属性，只会更新非零值的字段
		err = _self.db.Updates(model).Error
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
		tx = _self.db.Unscoped()
	} else {
		tx = _self.db
	}
	result := tx.Delete(req.Value, req.GetIds())
	if result.Error != nil {
		return ResultFailureErr(result.Error)
	}
	if result.RowsAffected <= 0 {
		_self.GetLogger().Error("删除失败,无删除记录")
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
