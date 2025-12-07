package main

import (
	"video-api/database"
	"video-api/handler"
	"video-api/pkg/config"
	"video-api/pkg/log"
	"video-api/pkg/ws"
	"video-api/repository"
	"video-api/router"
	"video-api/service"
)

func main() {
	config.Init()
	log.Init()
	defer log.Log.Sync()
	log.Log.Info("项目正在启动")
	database.InitDB()
	database.InitRedis()
	//启动webSocket管理器
	go ws.WebManager.Start()
	userRepo := repository.NewUserRepository(database.DB)
	videoRepo := repository.NewVideoRepository(database.DB)
	interactionRepo := repository.NewInteractionRepository(database.DB)
	socialRepo := repository.NewSocialRepository(database.DB)
	msgRepo := repository.NewMessageRepository(database.DB)

	userService := service.NewUserService(userRepo, socialRepo, database.RDB, database.Ctx)
	videoService := service.NewVideoService(userRepo, videoRepo, database.RDB, database.Ctx)
	interactionService := service.NewInteractionService(interactionRepo, userRepo, database.RDB, database.Ctx)
	socialService := service.NewSocialService(socialRepo, userRepo, database.RDB, database.Ctx)
	msgService := service.NewMessageService(msgRepo)

	userHandler := handler.NewUserHandler(userService)
	videoHandler := handler.NewVideoHandler(videoService)
	interactionHandler := handler.NewInteractionHandler(interactionService)
	socialHandler := handler.NewSocialHandler(socialService)
	chatHandler := handler.NewChatHandler(msgService)

	r := router.SetupRouter(userHandler, videoHandler, interactionHandler, socialHandler, chatHandler)
	r.Run(":" + config.Conf.Server.Port)
}
