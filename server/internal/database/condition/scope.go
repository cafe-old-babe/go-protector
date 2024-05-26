package condition

import (
	"fmt"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/custom/c_type"
	"gorm.io/gorm"
	"reflect"
)

type subSliceFunc func(db *gorm.DB, column string, slice interface{}) *gorm.DB

var subSliceCondMap map[c_type.SliceCondition]subSliceFunc

func init() {
	subSliceCondMap = make(map[c_type.SliceCondition]subSliceFunc)
	subSliceCondMap[consts.SliceIn] = func(db *gorm.DB, column string, slice interface{}) *gorm.DB {
		if slice == nil {
			return db
		}
		if reflect.Indirect(reflect.ValueOf(slice)).IsZero() {
			return db
		}
		return db.Or(fmt.Sprintf("%s %s ?", column, consts.SliceIn), slice)
	}
	subSliceCondMap[consts.SliceIn] = func(db *gorm.DB, column string, slice interface{}) *gorm.DB {
		if slice == nil {
			return db
		}
		if reflect.Indirect(reflect.ValueOf(slice)).IsZero() {
			return db
		}
		return db.Where(fmt.Sprintf("%s %s ?", column, consts.SliceNotIn), slice)

	}
}

// Paginate 分页查询
func Paginate(page base.IPagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page != nil {
			offset := (page.GetPageIndex() - 1) * page.GetPageSize()

			db = db.Offset(offset).Limit(page.GetPageSize())
		}

		return db
	}
}

// Like 封装like
func Like(column, arg string) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		if len(column) <= 0 || len(arg) <= 0 {
			return db
		}
		return db.Where(fmt.Sprintf("%s like ?", column), formatLike(arg))

	}

}

// LikeRight 封装like
func LikeRight(column, arg string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(column) <= 0 || len(arg) <= 0 {
			return db
		}
		return db.Where(fmt.Sprintf("%s like ?", column), formatLikeRight(arg))
	}

}

// LikeLeft 封装like
func LikeLeft(column, arg string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(column) <= 0 || len(arg) <= 0 {
			return db
		}
		return db.Where(fmt.Sprintf("%s like ?", column), formatLikeLeft(arg))
	}
}

func formatLike(arg string) string {
	return fmt.Sprintf("%%%s%%", arg)
}

func formatLikeRight(arg string) string {
	return fmt.Sprintf("%%%s", arg)
}
func formatLikeLeft(arg string) string {
	return fmt.Sprintf("%s%%", arg)
}

func EqStr(column, arg string) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		if len(column) <= 0 || len(arg) <= 0 {
			return db
		}
		return db.Where(column+" = ?", arg)
	}
}

// In 生成in条件
func In[T any](column string, slice []T) func(*gorm.DB) *gorm.DB {
	return generateSliceFunc(column, slice, consts.SliceIn)
}

// NotIn 生成not in条件
func NotIn[T any](column string, slice []T) func(*gorm.DB) *gorm.DB {
	return generateSliceFunc(column, slice, consts.SliceNotIn)
}

// generateSliceFunc 生成针对slice的方法
func generateSliceFunc[T any](column string, slice []T,
	condition c_type.SliceCondition) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if len(column) <= 0 || len(condition) <= 0 || len(slice) <= 0 {

			return db
		}

		sliceFunc, ok := subSliceCondMap[condition]
		if !ok {
			return db
		}
		start := 0
		sliceLen := len(slice)
		end := sliceLen
		batch := 500
		if end > batch {
			end = batch
		}

		subCond := db.Session(&gorm.Session{
			NewDB: true,
		})
		for {
			subCond = sliceFunc(subCond, column, slice[start:end])
			//subCond = subCond.Or(fmt.Sprintf("%s %s ?", column, condition), slice[start:end])

			if end >= sliceLen {
				break
			}
			start = end
			end += batch
			if end > sliceLen {
				end = sliceLen
			}
		}
		db.Where(subCond)
		return db
	}
}
