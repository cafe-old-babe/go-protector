package scope

import (
	"fmt"
	"go-protector/server/core/base"
	"gorm.io/gorm"
)

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
