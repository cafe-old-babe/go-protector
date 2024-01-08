package database

import (
	"context"
	"fmt"
	"go-protector/server/core/current"
	"go-protector/server/core/local"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"strings"
	"sync"
)

var _db *gorm.DB

var once sync.Once

func SetDB(db *gorm.DB) {
	once.Do(func() {
		databaseName := db.Migrator().CurrentDatabase()
		_ = db.Callback().Create().Before("gorm:create").Register("update_created_by", func(db *gorm.DB) {
			for _, field := range db.Statement.Schema.Fields {
				_ = field
			}
		})
		_ = db.Callback().Update().Before("gorm:update").Register("update_updated_by", func(db *gorm.DB) {

			if db.Error != nil {
				return
			}
			var (
			//selectColumns, restricted = db.Statement.SelectAndOmitColumns(false, true)
			)
			//db.Statement.SetColumn()
			if db.Statement.Schema == nil {

				set := callbacks.ConvertToAssignments(db.Statement)
				for _, assignment := range set {
					if strings.ToLower(assignment.Column.Name) == "update_by" {
						return
					}
				}
				//
				newDb := db.Session(&gorm.Session{NewDB: true})
				newDb.Logger.LogMode(logger.Warn)
				var count int64
				newDb.Table("information_schema.columns").
					Where("table_schema = ? and table_name = ? and column_name = ?",
						databaseName, db.Statement.Table, "updated_by").Count(&count)
				if count <= 0 {
					return
				}

				db.Statement.AddClause(append(set, clause.Assignment{
					Column: clause.Column{
						Name: "updated_by",
					},
					Value: current.GetUserId(db.Statement.Context),
					//Value: 2,
				}))

			} else {
				if field, ok := db.Statement.Schema.FieldsByName["UpdatedBy"]; ok {
					value, zero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
					fmt.Printf("field: %v, value: %v, zero: %v\n", field.Name, value, zero)
					_ = db.Statement.AddError(field.Set(db.Statement.Context, db.Statement.ReflectValue, 2))

				}
			}

		})

		_db = db
	})
}

func GetDB(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(local.CtxKeyDB).(*gorm.DB); ok {
		return db
	}

	return _db.WithContext(ctx)
}
