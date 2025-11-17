package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"video-api/model"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/video_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = DB.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	err = DB.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
