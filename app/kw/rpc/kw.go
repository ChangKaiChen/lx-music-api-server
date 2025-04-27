package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/kw"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
)

type KwHandler struct{}

func (s *KwHandler) KwMusicUrl(ctx context.Context, req *kw.KwRequest) (*kw.KwResponse, error) {
	log := logger.GetLogger()
	resp := &kw.KwResponse{Code: consts.ServerErrorCode}
	//songId := req.SongId
	//quality := req.Quality
	// TODO 获取url
	resp.Msg = "Not supported yet"
	log.Info("", resp.Msg)
	return resp, nil
}
