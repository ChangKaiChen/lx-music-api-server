package handler

import (
	"context"
	rpcAuth "github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc/auth"
	rpc "github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc/script"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
)

func Script(ctx context.Context, c *app.RequestContext) {
	log := logger.GetLogger()
	path := string(c.Path())
	url := strings.Split(c.Request.URI().String(), "/script")[0]
	key := c.Query("key")
	checkUpdate := c.Query("checkUpdate")
	respAuth, err := rpcAuth.AuthRPC(ctx, c, key)
	if err != nil {
		log.Errorf(path, "auth rpc err: %v", err)
		return
	}
	if !respAuth.Success {
		log.Log(path, key, "", "", "", "Unauthorized")
		response.ForbiddenResponse(c, "Unauthorized")
		return
	}
	respScript, err := rpc.ScriptRPC(ctx, c, url, key, checkUpdate)
	if err != nil {
		log.Errorf(path, "script rpc err: %v", err)
		return
	}
	if !respScript.Success {
		response.ServerErrorResponse(c, respScript.Res.Msg)
		return
	}
	if respScript.Res.Content != "" {
		response.SuccessResponse(c, respScript.Res.Content)
		return
	}
	response.SuccessResponse(c, respScript.Res)
}
