package middleware

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/opensergo/sentinel/adapter"
)

func SentinelMW() app.HandlerFunc {
	initSentinel()
	log := logger.GetLogger()
	return adapter.SentinelServerMiddleware(
		adapter.WithServerResourceExtractor(func(c context.Context, ctx *app.RequestContext) string {
			return "url"
		}),
		adapter.WithServerBlockFallback(func(ctx context.Context, c *app.RequestContext) {
			log.Errorf("", "frequent requests have been rejected by the gateway. clientIP: %v\n", c.ClientIP())
			c.AbortWithStatusJSON(consts.TooManyStatusCode, map[string]interface{}{
				"code": consts.TooManyStatusCode,
				"msg":  "Too many requests, please try again later",
			})
		}),
	)
}
func initSentinel() {
	log := logger.GetLogger()
	err := sentinel.InitDefault()
	if err != nil {
		log.Errorf("", "Unexpected error: %+v", err)
	}
	// limit QPS to 100
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "url",
			Threshold:              consts.SentinelThreshold,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       consts.SentinelStatIntervalInMs,
		},
	})
	if err != nil {
		log.Errorf("", "Unexpected error: %+v", err)
		return
	}
}
