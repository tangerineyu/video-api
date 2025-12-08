package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"video-api/pkg/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		status := c.Writer.Status()

		log.Log.Info("HTTP Request",
			zap.String("path", path),
			zap.String("query", query),
			zap.String("method", c.Request.Method),
			zap.Int("status", status),
			zap.Duration("cost", cost),
			zap.String("ip", c.ClientIP()),
		)
	}
}
