package handler

import (
	"context"
	"fmt"
	rpcAuth "github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc/auth"
	rpcKg "github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc/kg"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func KgMusicUrl(ctx context.Context, c *app.RequestContext) {
	log := logger.GetLogger()
	key := string(c.GetHeader("X-Request-Key"))
	path := string(c.Path())
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
	songId := c.Param("uid")
	quality := c.Param("quality")
	respKg, err := rpcKg.KgMusicUrlRPC(ctx, c, songId, quality)
	if err != nil {
		log.Errorf(path, "kgmusic url rpc err: %v", err)
		return
	}
	msg := fmt.Sprintf("target: %v, result: %v", string(c.Request.URI().FullURI()), respKg.Data)
	log.Log(path, key, "", "", "", msg)
	response.SuccessResponse(c, respKg)
}
