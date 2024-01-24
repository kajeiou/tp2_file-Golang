package db

import (
	"log"
	"tp2/dictionary"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitializeDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("db/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	DB.AutoMigrate(&dictionary.Word{})
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get the underlying SQL DB:", err)
		return
	}
	sqlDB.Close()
}
