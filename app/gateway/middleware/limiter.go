package middleware

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/limiter"
	"github.com/cloudwego/hertz/pkg/app"
)

func LimiterMW() app.HandlerFunc {
	l := limiter.NewLimiter(
		limiter.WithRate(consts.SingleKeyRequestRate),
		limiter.WithWindow(consts.SingleKeyRequestWindow),
	)
	go l.StartCleanup()
	return func(ctx context.Context, c *app.RequestContext) {
		key := l.GetKey(ctx, c)
		if !l.IsAllowed(key) {
			c.JSON(consts.TooManyStatusCode, map[string]any{
				"code": consts.TooManyCode,
				"msg":  "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
