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
	interactionRepo := repository.NewInteractionRepository(database.DB)
	socialRepo := repository.NewSocialRepository(database.DB)

	userService := service.NewUserService(userRepo, database.RDB, database.Ctx)
	videoService := service.NewVideoService(userRepo, videoRepo, database.RDB, database.Ctx)
	interactionService := service.NewInteractionService(interactionRepo, userRepo, database.RDB, database.Ctx)
	socialService := service.NewSocialService(socialRepo, userRepo, database.RDB, database.Ctx)

	userHandler := handler.NewUserHandler(userService)
	videoHandler := handler.NewVideoHandler(videoService)
	interactionHander := handler.NewInteractionHandler(interactionService)
	socialHandler := handler.NewSocialHandler(socialService)

	r := router.SetupRouter(userHandler, videoHandler, interactionHander, socialHandler)
	r.Run(":8080")
}
