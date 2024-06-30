package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}
