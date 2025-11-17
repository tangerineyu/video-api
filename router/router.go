package router

import (
	"github.com/gin-gonic/gin"
	"video-api/handler"
	"video-api/middleware"
)

func SetupRouter(
	userHandler *handler.UserHandler,
) *gin.Engine {
	r := gin.Default()
	apigroup := r.Group("/")
	userGroup := apigroup.Group("/user")
	{

	}
	return r
}
