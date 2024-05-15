package database

import (
	"context"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"strings"
	"sync"
)

var _db *gorm.DB

var once sync.Once
var databaseName string

func SetDB(db *gorm.DB) {
	once.Do(func() {
		databaseName = db.Migrator().CurrentDatabase()
		_ = db.Callback().Create().Before("gorm:create").Register("created_by", func(db *gorm.DB) {
			if db.Error != nil {
				return
			}
			userId := current.GetUserId(db.Statement.Context)
			if userId <= 0 {
				return
			}
			if db.Statement.Schema == nil {
				// todo 留着
			} else {
				if field, ok := db.Statement.Schema.FieldsByName["CreatedBy"]; ok {
					_, zero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
					//fmt.Printf("field: %v, value: %v, zero: %v\n", field.Name, value, zero)
					if zero {
						_ = db.Statement.AddError(field.Set(db.Statement.Context, db.Statement.ReflectValue, userId))
					}
				}
			}
		})
		_ = db.Callback().Update().Before("gorm:update").Register("updated_by", updateByFunc)
		_db = db
	})
}

func GetDB(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(consts.CtxKeyDB).(*gorm.DB); ok {
		return db
	}

	return _db.WithContext(ctx)
}

func updateByFunc(db *gorm.DB) {

	if db.Error != nil {
		return
	}

	var (
		userId = current.GetUserId(db.Statement.Context)
		//selectColumns, restricted = db.Statement.SelectAndOmitColumns(false, true)
	)
	if userId <= 0 {
		return
	}
	set := callbacks.ConvertToAssignments(db.Statement)

	//db.Statement.SetColumn()
	if db.Statement.Schema == nil {

		for _, assignment := range set {
			if strings.ToLower(assignment.Column.Name) == "updated_by" {
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
			Value: userId,
			//Value: 2,
		}))

	} else {
		if field, ok := db.Statement.Schema.FieldsByName["UpdatedBy"]; ok {
			_, zero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
			//fmt.Printf("field: %v, value: %v, zero: %v\n", field.Name, value, zero)
			if !zero {
				return
			}
			db.Statement.AddClause(append(set, clause.Assignment{
				Column: clause.Column{
					Name: field.DBName,
				},
				Value: userId,
			}))
			//_ = db.Statement.AddError(field.Set(db.Statement.Context, db.Statement.ReflectValue, userId))

		}
	}

}
