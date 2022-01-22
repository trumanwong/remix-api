package middlewares

import (
	"remix-api/configs"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/trumanwong/go-internal/util"
	"net/http"
)

type AuthCheck struct {
}

func (this *AuthCheck) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if len(tokenString) < 7 {
			ctx.Abort()
			util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
			return
		}
		tokenString = tokenString[7:] //"Bearer "

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(configs.Config.Other.JwtKey), nil
		})

		if err != nil {
			util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if _, ok := claims["uuid"]; !ok {
				util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
				ctx.Abort()
			}

			ctx.Set("uuid", claims["uuid"])
			ctx.Set("type", claims["type"])
		} else {
			util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
			ctx.Abort()
		}
		ctx.Next()
	}
}
