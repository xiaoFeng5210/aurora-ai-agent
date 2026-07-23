package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ApiTimers(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()
	elapsed := time.Since(start)
	// 如果elapsed大于1s的 打印请求耗时
	if elapsed > time.Second && elapsed < time.Second*2 {
		logger.Warn("request took long", zap.String("method", ctx.Request.Method), zap.Int("status", ctx.Writer.Status()), zap.String("path", ctx.Request.URL.Path), zap.Duration("elapsed", elapsed))
	}

	if elapsed >= time.Second*2 {
		logger.Warn("request took too long", zap.String("method", ctx.Request.Method), zap.Int("status", ctx.Writer.Status()), zap.String("path", ctx.Request.URL.Path), zap.Duration("elapsed", elapsed))
	}
}
