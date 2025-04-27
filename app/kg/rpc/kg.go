package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/kg"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
)

var KgQualityHashMap = map[string]string{
	"128k":      "hash_128",
	"320k":      "hash_320",
	"flac":      "hash_flac",
	"flac24bit": "hash_high",
}
var KgQualityMap = map[string]string{
	"128k":      "128",
	"320k":      "320",
	"flac":      "flac",
	"flac24bit": "high",
}

type KgHandler struct{}

func (s *KgHandler) KgMusicUrl(ctx context.Context, req *kg.KgRequest) (*kg.KgResponse, error) {
	log := logger.GetLogger()
	resp := &kg.KgResponse{Code: consts.ServerErrorCode}
	//songId := strings.ToLower(req.SongId)
	//quality := req.Quality
	// TODO 获取url
	resp.Msg = "Not supported yet"
	log.Info("", resp.Msg)
	return resp, nil
}
