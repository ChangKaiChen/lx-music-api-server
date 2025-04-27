package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/mg"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func MgMusicUrlRPC(ctx context.Context, c *app.RequestContext, songId, quality string) (*mg.MgResponse, error) {
	log := logger.GetLogger()
	resp := &mg.MgResponse{}
	req := &mg.MgRequest{
		SongId:  songId,
		Quality: quality,
	}
	resp, err := (*rpc.MgClient).MgMusicUrl(ctx, req)
	if err != nil {
		response.SuccessResponse(c, resp.Msg)
		log.Errorf("", "MgMusicUrlRPC call failed: %v", err)
		return resp, err
	}
	return resp, nil
}
