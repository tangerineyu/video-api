package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"video-api/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	//从环境变量中获取DB地址，默认为127.0.0.1
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "123456"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "video_db"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println("dsn:", dsn)
	//
	//dsn := "root:123456@tcp(127.0.0.1:3306)/video_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	var dbConn *gorm.DB
	for i := 0; i < 120; i++ {
		dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err == nil {
			sqlDB, _ := dbConn.DB()
			if err = sqlDB.Ping(); err == nil {
				log.Println("successfully connected to database")
				DB = dbConn
				break
			}
		}
		//fmt.Println("数据库还没准备好，重试中(%d)...\n", i+1)
		time.Sleep(1 * time.Second)
	}
	//DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if DB == nil {
		log.Fatal("数据库连接失败，退出程序，最后一次错误: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 自动迁移模型
	err = DB.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.UserFavorite{},
		&model.UserRelation{},
		&model.Message{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
