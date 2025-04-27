package handler

import (
	"context"
	"fmt"
	rpcAuth "github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc/auth"
	rpcTx "github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc/tx"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func TxMusicUrl(ctx context.Context, c *app.RequestContext) {
	log := logger.GetLogger()
	path := string(c.Path())
	key := string(c.GetHeader("X-Request-Key"))
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
	respTx, err := rpcTx.TxMusicUrlRPC(ctx, c, songId, quality)
	if err != nil {
		log.Errorf(path, "txmusic url rpc err: %v", err)
		return
	}
	msg := fmt.Sprintf("target: %v, result: %v", string(c.Request.URI().FullURI()), respTx.Data)
	log.Log(path, key, "", "", "", msg)
	response.SuccessResponse(c, respTx)
}
