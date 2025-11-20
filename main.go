package main

import (
	"video-api/database"
	"video-api/handler"
	"video-api/repository"
	"video-api/router"
	"video-api/service"
)

func main() {
	database.InitDB()
	database.InitRedis()

	userRepo := repository.NewUserRepository(database.DB)
	videoRepo := repository.NewVideoRepository(database.DB)

	userService := service.NewUserService(userRepo, database.RDB, database.Ctx)
	videoService := service.NewVideoService(userRepo, videoRepo, database.RDB, database.Ctx)

	userHandler := handler.NewUserHandler(userService)
	videoHandler := handler.NewVideoHandler(videoService)

	r := router.SetupRouter(userHandler, videoHandler)
	r.Run(":8080")
}
