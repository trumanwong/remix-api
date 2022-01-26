package middlewares

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/trumanwong/go-internal/util"
	"net/http"
	"time"
)

type Throttle struct {
	lmt *limiter.Limiter
}

func NewThrottle(max float64) *Throttle {
	lmt := tollbooth.NewLimiter(max, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
	})
	return &Throttle{lmt: lmt}
}

func (this *Throttle) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		this.lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).SetMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
		err := tollbooth.LimitByRequest(this.lmt, ctx.Writer, ctx.Request)
		if err != nil {
			util.Response(ctx, nil, http.StatusTooManyRequests, "Too Many Attempts.")
			ctx.Abort()
		}
		ctx.Next()
	}
}
