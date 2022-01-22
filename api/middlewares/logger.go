package middlewares

import (
	"remix-api/internal/logs"
	"github.com/gin-gonic/gin"
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
			// 请求头
			Header: ctx.Request.Header,
			// 请求方式
			Method: ctx.Request.Method,
			// 响应状态码
			StatusCode: ctx.Writer.Status(),
			// 执行时间
			ExecuteTime: endTime.Sub(startTime) / time.Millisecond,
			CreatedAt:   time.Now(),
		})
	}
}
