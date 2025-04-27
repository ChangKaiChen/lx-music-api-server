package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/wy"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func WyMusicUrlRPC(ctx context.Context, c *app.RequestContext, songId, quality string) (*wy.WyResponse, error) {
	log := logger.GetLogger()
	resp := &wy.WyResponse{}
	req := &wy.WyRequest{
		SongId:  songId,
		Quality: quality,
	}
	resp, err := (*rpc.WyClient).WyMusicUrl(ctx, req)
	if err != nil {
		response.SuccessResponse(c, resp.Msg)
		log.Errorf("", "WyMusicUrlRPC call failed: %v", err)
		return resp, err
	}
	return resp, nil
}
