package response

import (
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/cloudwego/hertz/pkg/app"
)

func ServerErrorResponse(c *app.RequestContext, msg string) {
	c.JSON(consts.ServerErrorStatusCode, map[string]any{
		"code": consts.ServerErrorCode,
		"msg":  msg,
	})
}
func ForbiddenResponse(c *app.RequestContext, msg string) {
	c.JSON(consts.ForbiddenStatusCode, map[string]any{
		"code": consts.ForbiddenCode,
		"msg":  msg,
	})
}
func SuccessResponse(c *app.RequestContext, data any) {
	c.JSON(consts.SuccessStatusCode, data)
}
