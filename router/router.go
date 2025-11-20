package router

import (
	"video-api/handler"
	"video-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	videoHandler *handler.VideoHandler,
) *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./uploads")
	apigroup := r.Group("/")
	userGroup := apigroup.Group("/user")
	{
		userGroup.POST("register/", userHandler.Register)
		userGroup.POST("login/", userHandler.Login)
		userGroup.GET("/", middleware.AuthMiddleware(), userHandler.GetUserInfo)
		userGroup.POST("avatar/upload/", middleware.AuthMiddleware(), userHandler.UploadAvatar)
	}
	feedGroup := r.Group("/feed")
	{
		feedGroup.GET("/", videoHandler.Feed)
	}
	publishGroup := r.Group("/publish")
	{
		publishGroup.POST("action/", middleware.AuthMiddleware(), videoHandler.Publish)
		publishGroup.GET("list/", middleware.AuthMiddleware(), videoHandler.List)
	}

	return r
}
