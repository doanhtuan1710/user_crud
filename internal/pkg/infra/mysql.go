package infra

import (
	"user_crud/internal/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewMySQLDB() (db *gorm.DB, err error) {

	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	// Auto migrate db
	db.AutoMigrate(&entity.User{})

	return
}
