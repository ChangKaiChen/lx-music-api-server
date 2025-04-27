package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/mg"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
)

type MgHandler struct{}

func (s *MgHandler) MgMusicUrl(ctx context.Context, req *mg.MgRequest) (*mg.MgResponse, error) {
	log := logger.GetLogger()
	resp := &mg.MgResponse{Code: consts.ServerErrorCode}
	//songId := req.SongId
	//quality := req.Quality
	// TODO 获取url
	resp.Msg = "Not supported yet"
	log.Info("", resp.Msg)
	return resp, nil
}
