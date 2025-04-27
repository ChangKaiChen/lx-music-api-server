package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/tx"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func TxMusicUrlRPC(ctx context.Context, c *app.RequestContext, songId, quality string) (*tx.TxResponse, error) {
	log := logger.GetLogger()
	resp := &tx.TxResponse{}
	req := &tx.TxRequest{
		SongId:  songId,
		Quality: quality,
	}
	resp, err := (*rpc.TxClient).TxMusicUrl(ctx, req)
	if err != nil {
		response.SuccessResponse(c, resp.Msg)
		log.Errorf("", "TxMusicUrlRPC call failed: %v", err)
		return resp, err
	}
	return resp, nil
}
