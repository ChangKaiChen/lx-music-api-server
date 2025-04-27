package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChangKaiChen/lx-music-api-server/app/tx/config"
	"github.com/ChangKaiChen/lx-music-api-server/app/tx/crypto"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/tx"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/cache"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"time"
)

var fileInfo = map[string]map[string]string{
	"128k": {
		"e": ".mp3",
		"h": "M500",
	},
	"320k": {
		"e": ".mp3",
		"h": "M800",
	},
	"flac": {
		"e": ".flac",
		"h": "F000",
	},
	"flac24bit": {
		"e": ".flac",
		"h": "RS01",
	},
	"dolby": {
		"e": ".flac",
		"h": "Q000",
	},
	"master": {
		"e": ".flac",
		"h": "AI00",
	},
}
var TxQualityReverseMap = map[string]string{
	"M500": "128k",
	"M800": "320k",
	"F000": "flac",
	"RS01": "flac24bit",
	"Q000": "dolby",
	"AI00": "master",
}

type Comm struct {
	Ct  string `json:"ct"`
	Cv  string `json:"cv"`
	Uin string `json:"uin"`
}

type Param struct {
	SongType int    `json:"song_type"`
	SongMid  string `json:"song_mid"`
}

type InfoBodyReq struct {
	Module string `json:"module"`
	Method string `json:"method"`
	Param  Param  `json:"param"`
}

type InfoBody struct {
	Comm        Comm        `json:"comm"`
	InfoBodyReq InfoBodyReq `json:"req"`
}
type RequestBody struct {
	RequestReq  RequestReq  `json:"req"`
	RequestComm RequestComm `json:"comm"`
}
type RequestReq struct {
	Module       string       `json:"module"`
	Method       string       `json:"method"`
	RequestParam RequestParam `json:"param"`
}
type RequestParam struct {
	Filename  []string `json:"filename"`
	Guid      string   `json:"guid"`
	SongMid   []string `json:"songmid"`
	SongType  []int    `json:"songtype"`
	Uin       string   `json:"uin"`
	LoginFlag int      `json:"loginflag"`
	Platform  string   `json:"platform"`
}
type RequestComm struct {
	QQ     string `json:"qq"`
	AuthST string `json:"authst"`
	CT     string `json:"ct"`
	CV     string `json:"cv"`
	V      string `json:"v"`
}
type InfoResponse struct {
	InfoReq InfoReq `json:"req"`
}
type InfoReq struct {
	InfoData InfoData `json:"data"`
}
type InfoData struct {
	TrackInfo TrackInfo `json:"track_info"`
}
type TrackInfo struct {
	File File `json:"file"`
}
type File struct {
	MediaMid string `json:"media_mid"`
}

type Response struct {
	Code int `json:"code"`
	Req  Req `json:"req"`
}

type Req struct {
	Data Data `json:"data"`
}
type Data struct {
	MidUrlInfo []MidUrlInfo `json:"midurlinfo"`
}
type MidUrlInfo struct {
	Purl     string `json:"purl"`
	Filename string `json:"filename"`
}
type TxHandler struct{}

func (s *TxHandler) TxMusicUrl(ctx context.Context, req *tx.TxRequest) (*tx.TxResponse, error) {
	log := logger.GetLogger()
	resp := &tx.TxResponse{Code: consts.ServerErrorCode}
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
	key := fmt.Sprintf("%s-%s-%s", consts.TxServiceName, songId, quality)
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
	if user.Uin == "" || user.QQMusicKey == "" {
		resp.Msg = "User is none"
		log.Info("", resp.Msg)
		return resp, nil
	}
	url := "https://u.y.qq.com/cgi-bin/musics.fcg?format=json&sign="
	infoBody := InfoBody{
		Comm: Comm{
			Ct:  "19",
			Cv:  "1859",
			Uin: "0",
		},
		InfoBodyReq: InfoBodyReq{
			Module: "music.pf_song_detail_svr",
			Method: "get_song_detail_yqq",
			Param: Param{
				SongType: 0,
				SongMid:  songId,
			},
		},
	}
	jsonInfoBody, err := json.Marshal(infoBody)
	if err != nil {
		resp.Msg = "JSON encoding error"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	jsonStr := string(jsonInfoBody)
	//jsonStr = strings.ReplaceAll(jsonStr, ":", ": ")
	//jsonStr = strings.ReplaceAll(jsonStr, ",", ", ")
	sign, err := crypto.Sign(jsonStr, false)
	if err != nil {
		resp.Msg = "sign encrypt error"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	r1, err := http.NewRequest("POST", url+sign, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		resp.Msg = "Server NewRequest failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(r1)
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
	infoRes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Msg = "Server IO read failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	var info InfoResponse
	err = json.Unmarshal(infoRes, &info)
	if err != nil {
		resp.Msg = "Server JSON decode failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	strMediaMid := info.InfoReq.InfoData.TrackInfo.File.MediaMid
	requestBody := RequestBody{
		RequestReq: RequestReq{
			Module: "music.vkey.GetVkey",
			Method: "UrlGetVkey",
			RequestParam: RequestParam{
				Filename:  []string{fileInfo[quality]["h"] + strMediaMid + fileInfo[quality]["e"]},
				Guid:      "114514",
				SongMid:   []string{songId},
				SongType:  []int{0},
				Uin:       user.Uin,
				LoginFlag: 1,
				Platform:  "20",
			},
		},
		RequestComm: RequestComm{
			QQ:     user.Uin,
			AuthST: user.QQMusicKey,
			CT:     "26",
			CV:     "2010101",
			V:      "2010101",
		},
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		resp.Msg = "JSON encoding error"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	jsonStr = string(jsonBody)
	//jsonStr = strings.ReplaceAll(jsonStr, ":", ": ")
	//jsonStr = strings.ReplaceAll(jsonStr, ",", ", ")
	sign, err = crypto.Sign(jsonStr, true)
	if err != nil {
		resp.Msg = "sign encrypt error"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	r2, err := http.NewRequest("POST", url+sign, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		resp.Msg = "Server NewRequest failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	httpResp, err = httpClient.Do(r2)
	if err != nil {
		resp.Msg = "Server request failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	response, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Msg = "Server IO read failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	var res Response
	err = json.Unmarshal(response, &res)
	if err != nil {
		resp.Msg = "Server JSON decode failed"
		log.Errorf("", resp.Msg+": %v", err)
		return resp, err
	}
	if res.Code != 0 || res.Req.Data.MidUrlInfo[0].Purl == "" {
		resp.Code = consts.NotFoundCode
		resp.Msg = "qqmusic_key invalid"
		log.Error("", resp.Msg)
		return resp, errors.New("qqmusic_key invalid")
	}
	data := res.Req.Data.MidUrlInfo[0]
	resp.Code = consts.SuccessCode
	resp.Msg = "Success"
	resp.Data = "http://ws.stream.qqmusic.qq.com/" + data.Purl
	resp.Extra = &tx.Extra{
		Cache: false,
		Quality: &tx.Quality{
			Target:  quality,
			Result_: TxQualityReverseMap[strings.Split(data.Filename, ".")[0][:4]],
		},
		Expire: &tx.Expire{
			Time:      int64(config.GetConf().ExpireTime) + time.Now().Unix() - 120,
			CanExpire: true,
		},
	}
	rdb.Set(key, resp, config.GetConf().ExpireTime)
	log.Log("", "", songId, quality, resp.Data, "")
	return resp, nil
}
