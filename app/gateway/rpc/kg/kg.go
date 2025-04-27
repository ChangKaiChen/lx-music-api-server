package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/kg"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func KgMusicUrlRPC(ctx context.Context, c *app.RequestContext, songId, quality string) (*kg.KgResponse, error) {
	log := logger.GetLogger()
	resp := &kg.KgResponse{}
	req := &kg.KgRequest{
		SongId:  songId,
		Quality: quality,
	}
	resp, err := (*rpc.KgClient).KgMusicUrl(ctx, req)
	if err != nil {
		response.SuccessResponse(c, resp.Msg)
		log.Errorf("", "KgMusicUrlRPC call failed: %v", err)
		return resp, err
	}
	return resp, nil
}
