package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger 日志中间件（记录请求信息）
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 请求信息
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// 记录日志
		log.Info().
			Str("method", reqMethod).
			Str("uri", reqURI).
			Int("status", statusCode).
			Str("ip", clientIP).
			Dur("latency", latency).
			Msg("请求处理完成")
	}
}
