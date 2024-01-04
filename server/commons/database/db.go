package database

import (
	"context"
	"go-protector/server/commons/local"
	"gorm.io/gorm"
	"sync"
)

var _db *gorm.DB

var once sync.Once

func SetDB(db *gorm.DB) {
	once.Do(func() {
		_db = db
	})
}

func GetDB(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(local.CtxKeyDB).(*gorm.DB); ok {
		return db
	}

	return _db.WithContext(ctx)
}
