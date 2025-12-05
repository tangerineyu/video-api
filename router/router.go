package router

import (
	"video-api/handler"
	"video-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	videoHandler *handler.VideoHandler,
	interactionHandler *handler.InteractionHandler,
	socialHandler *handler.SocialHandler,
	chatHandler *handler.ChatHandler,
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
		userGroup.GET("mfa/qrcode", middleware.AuthMiddleware(), userHandler.GenerateMFA)
		userGroup.POST("mfa/bind", middleware.AuthMiddleware(), userHandler.BindMFA)
		//userGroup.POST("search/image", userHandler.ImageSearch)
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
	favoriteGroup := apigroup.Group("favorite")
	{
		favoriteGroup.POST("action/", middleware.AuthMiddleware(), interactionHandler.FavoriteAction)
		favoriteGroup.GET("list/", middleware.AuthMiddleware(), interactionHandler.FavoriteList)
	}
	commentGroup := r.Group("comment")
	{
		commentGroup.POST("action/", middleware.AuthMiddleware(), interactionHandler.CommentAction)
		commentGroup.GET("list/", interactionHandler.CommentList)
	}
	relationGroup := apigroup.Group("relation")
	{
		relationGroup.POST("action/", middleware.AuthMiddleware(), socialHandler.RelationAction)
		relationGroup.GET("follow/list/", socialHandler.FollowList)
		relationGroup.GET("follower/list/", socialHandler.FollowerList)
		relationGroup.GET("friend/list/", middleware.AuthMiddleware(), socialHandler.FriendList)
	}
	chatGroup := r.Group("chat")
	{
		chatGroup.GET("ws", chatHandler.Connect)
		chatGroup.GET("history", middleware.AuthMiddleware(), chatHandler.GetHistory)
	}
	// 模拟点击
	apigroup.POST("video/visit/:id", videoHandler.VisitVideo)
	// 排行榜
	apigroup.GET("rank/popular", videoHandler.PopularRank)
	return r
}
