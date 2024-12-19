package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}

func MigrateModel(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
