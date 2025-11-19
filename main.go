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

	userService := service.NewUserService(userRepo, database.RDB, database.Ctx)

	userHandler := handler.NewUserHandler(userService)

	r := router.SetupRouter(userHandler)
	r.Run(":8080")
}
