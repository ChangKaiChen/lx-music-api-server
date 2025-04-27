package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChangKaiChen/lx-music-api-server/app/wy/config"
	"github.com/ChangKaiChen/lx-music-api-server/app/wy/crypto"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/wy"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/cache"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

var WyQualityMap = map[string]string{
	"128k":      "standard",
	"192k":      "higher",
	"320k":      "exhigh",
	"flac":      "lossless",
	"flac24bit": "hires",
}
var WyQualityReverseMap = map[string]string{
	"standard": "128k",
	"higher":   "192k",
	"exhigh":   "320k",
	"lossless": "flac",
	"hires":    "flac24bit",
}

type Response struct {
	Data []Data `json:"data"`
}
type Data struct {
	Url   string `json:"url"`
	Level string `json:"level"`
}
type ReqBody struct {
	IDS        string `json:"ids"`
	Level      string `json:"level"`
	EncodeType string `json:"encodeType"`
}

type WyHandler struct{}

func (s *WyHandler) WyMusicUrl(ctx context.Context, req *wy.WyRequest) (*wy.WyResponse, error) {
	log := logger.GetLogger()
	resp := &wy.WyResponse{Code: consts.ServerErrorCode}
	songId := req.SongId
	quality := req.Quality
	exists := slices.Contains(config.GetConf().Quality, quality)
	if !exists {
		resp.Code = consts.ForbiddenCode
		resp.Msg = "Do not support this quality"
		log.Info("", resp.Msg)
		return resp, nil
	}
	rdb := cache.GetCache()
	key := fmt.Sprintf("%s-%s-%s", consts.WyServiceName, songId, quality)
	err := rdb.Get(key, resp)
	if err == nil {
		resp.Extra.Cache = true
		return resp, nil
	}
	if len(config.GetConf().Users) == 0 {
		resp.Msg = "Users are empty"
		log.Info("", resp.Msg)
		return resp, nil
	}
	user := config.GetConf().Users[rand.Intn(len(config.GetConf().Users))]
	if user.Cookie == "" {
		resp.Msg = "User is none"
		log.Info("", resp.Msg)
		return resp, nil
	}
	path := "/api/song/enhance/player/url/v1"
	requestUrl := "https://interface.music.163.com/eapi/song/enhance/player/url/v1"
	ids, err := json.Marshal([]string{songId})
	if err != nil {
		resp.Msg = "JSON encoding error"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	requestBody := ReqBody{
		IDS:        string(ids),
		Level:      WyQualityMap[quality],
		EncodeType: "flac",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		resp.Msg = "JSON encoding error"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	formMap := crypto.EApiEncrypt(path, string(jsonBody))
	form := url.Values{}
	form.Add("params", formMap["params"])
	r, err := http.NewRequest("POST", requestUrl, bytes.NewBufferString(form.Encode()))
	if err != nil {
		resp.Msg = "Server NewRequest failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	r.Header.Set("Cookie", user.Cookie)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(r)
	if err != nil {
		resp.Msg = "Server request failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	defer func() {
		err = httpResp.Body.Close()
		if err != nil {
		}
	}()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Msg = "Server IO read failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		resp.Msg = "Server JSON decode failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	data := res.Data
	if len(data) == 0 || data[0].Url == "" {
		resp.Code = consts.NotFoundCode
		resp.Msg = "Request failed"
		log.Errorf("", resp.Msg+": %v", "response is empty")
		return resp, errors.New("server response is empty")
	}
	if data[0].Level != WyQualityMap[quality] {
		resp.Code = consts.NotFoundCode
		resp.Msg = fmt.Sprintf("Reject unmatched quality: expected=%s, got=%s", WyQualityMap[quality], data[0].Level)
		log.Info("", resp.Msg)
		return resp, nil
	}
	resp.Code = consts.SuccessCode
	resp.Msg = "Success"
	resp.Data = strings.Split(data[0].Url, "?")[0]
	resp.Extra = &wy.Extra{
		Cache: false,
		Quality: &wy.Quality{
			Target:  quality,
			Result_: WyQualityReverseMap[data[0].Level],
		},
		Expire: &wy.Expire{
			Time:      int64(config.GetConf().ExpireTime) + time.Now().Unix() - 120,
			CanExpire: true,
		},
	}
	rdb.Set(key, resp, config.GetConf().ExpireTime)
	log.Log("", "", songId, quality, resp.Data, "")
	return resp, nil
}
