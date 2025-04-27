package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/kw"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func KwMusicUrlRPC(ctx context.Context, c *app.RequestContext, songId, quality string) (*kw.KwResponse, error) {
	log := logger.GetLogger()
	resp := &kw.KwResponse{}
	req := &kw.KwRequest{
		SongId:  songId,
		Quality: quality,
	}
	resp, err := (*rpc.KwClient).KwMusicUrl(ctx, req)
	if err != nil {
		response.SuccessResponse(c, resp.Msg)
		log.Errorf("", "KwMusicUrlRPC call failed: %v", err)
		return resp, err
	}
	return resp, nil
}
