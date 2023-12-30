package database

import (
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

func GetDB() *gorm.DB {
	return _db
}
