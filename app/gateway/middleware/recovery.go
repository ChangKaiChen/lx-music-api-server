package middleware

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
)

func RecoveryMW() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
		log := logger.GetLogger()
		path := string(c.Path())
		log.Errorf(path, "[Recovery] InternalServiceError err=%v", err)
		c.JSON(consts.ServerErrorStatusCode, map[string]interface{}{
			"code": consts.ServerErrorCode,
			"msg":  "Internal service error, please try again later",
		})
	}))
}
