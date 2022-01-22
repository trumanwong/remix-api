package middlewares

import (
	"github.com/gin-gonic/gin"
	"remix-api/internal/logs"
	"time"
)

type Logger struct {
}

func (this *Logger) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		logs.Write(logs.AccessLog{
			// 请求路由
			Uri: ctx.Request.RequestURI,
			// 请求ip
			ClientIp: ctx.ClientIP(),
			// 请求方式
			Method: ctx.Request.Method,
			// 请求头
			Header: ctx.Request.Header,
			// 返回code
			StatusCode: ctx.Writer.Status(),
			// 执行时间
			ExecuteTime: endTime.Sub(startTime) / time.Millisecond,
			CreatedAt:   time.Now(),
		})
	}
}