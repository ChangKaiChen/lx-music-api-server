package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/script"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func ScriptRPC(ctx context.Context, c *app.RequestContext, url, key, checkUpdate string) (*script.ScriptResponse, error) {
	log := logger.GetLogger()
	resp := &script.ScriptResponse{}
	reqScript := &script.ScriptRequest{
		Key:         key,
		Url:         url,
		CheckUpdate: checkUpdate,
	}
	resp, err := (*rpc.ScriptClient).Script(ctx, reqScript)
	if err != nil {
		response.SuccessResponse(c, resp.Res)
		log.Errorf("", "ScriptRPC call failed: %v", err)
		return resp, err
	}
	return resp, nil
}
