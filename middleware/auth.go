package middleware

import (
	"errors"
	"fmt"
	"strings"
	"video-api/handler"
	"video-api/pkg/errno"
	"video-api/pkg/log"
	"video-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// 验证access token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("进入了中间件")
		authHeader := c.GetHeader("Authorization")
		fmt.Println("收到Header:", authHeader)
		if authHeader == "" {
			authHeader = c.Query("token")
			if authHeader == "" {
				log.Log.Warn("未提供token",
					zap.String("token", c.Query("token")))
				handler.Error(c, errno.ParamErr)
				c.Abort()
				return
			}
		}
		parts := strings.Split(authHeader, " ")
		tokenString := ""
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		} else {
			tokenString = authHeader
		}
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			fmt.Println("3.token解析失败，错误信息：", err)
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Log.Warn("token已经过期",
					zap.String("token", c.Query("token")),
					zap.Error(err))
				handler.Error(c, errno.TokenInvalidErr)

			} else {
				log.Log.Error("token验证失败",
					zap.String("token", c.Query("token")),
					zap.Error(err))
				handler.Error(c, errno.ServiceErr)
			}
			c.Abort()
			return
		}
		fmt.Println("4.token解析成功，用户ID：", claims.UserID)
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()

	}
}
