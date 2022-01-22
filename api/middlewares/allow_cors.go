package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"remix-api/configs"
	"strings"
)

type AllowCors struct {
}

func (this *AllowCors) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		allowOrigins := strings.Split(configs.Config.Other.AllowOrigins, ",")
		for _, v := range allowOrigins {
			if origin == v {
				ctx.Header("Access-Control-Allow-Origin", origin)
			}
		}
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,Accept,Authorization,X-Requested-With,X-XSRF-TOKEN,x-csrf-token,Cache-Control,token")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}

		ctx.Next()
	}
}
