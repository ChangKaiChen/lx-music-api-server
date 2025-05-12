package refresh

import (
	"bytes"
	"encoding/json"
	"github.com/ChangKaiChen/lx-music-api-server/app/tx/crypto"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"io"
	"net/http"
)

type ReqBody struct {
	ReqComm ReqComm `json:"comm"`
	ReqReq1 ReqReq1 `json:"req1"`
}
type ReqComm struct {
	FPersonality string `json:"fPersonality"`
	TmeLoginType string `json:"tmeLoginType"`
	QQ           string `json:"qq"`
	AuthST       string `json:"authst"`
	CT           string `json:"ct"`
	CV           string `json:"cv"`
	V            string `json:"v"`
	TmeAppID     string `json:"tmeAppID"`
}
type ReqReq1 struct {
	Module   string   `json:"module"`
	Method   string   `json:"method"`
	ReqParam ReqParam `json:"param"`
}
type ReqParam struct {
	StrMusicID string `json:"str_musicid"`
	MusicKey   string `json:"musickey"`
	RefreshKey string `json:"refresh_key"`
}
type Res struct {
	Code int  `json:"code"`
	Req1 Req1 `json:"req1"`
}
type Req1 struct {
	Data Data `json:"data"`
}
type Data struct {
	MusicKey string `json:"musickey"`
}

func QQMusicKeyRefresh(uin, musicKey string) string {
	log := logger.GetLogger()
	url := "https://u.y.qq.com/cgi-bin/musics.fcg?sign="
	body := ReqBody{
		ReqComm: ReqComm{
			FPersonality: "0",
			TmeLoginType: "2",
			QQ:           uin,
			AuthST:       musicKey,
			CT:           "11",
			CV:           "12080008",
			V:            "12080008",
			TmeAppID:     "qqmusic",
		},
		ReqReq1: ReqReq1{
			Module: "music.login.LoginServer",
			Method: "Login",
			ReqParam: ReqParam{
				StrMusicID: uin,
				MusicKey:   musicKey,
				RefreshKey: "",
			},
		},
	}
	if musicKey[0] != 'Q' {
		body.ReqComm.TmeLoginType = "1"
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return ""
	}
	jsonStr := string(jsonBody)
	sign, err := crypto.Sign(jsonStr, false)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return ""
	}
	r, err := http.NewRequest("POST", url+sign, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return ""
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return ""
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
		}
	}()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return ""
	}
	var res Res
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return ""
	}
	if res.Code != 0 {
		log.Error("", "musicKey refresh failed")
		return ""
	}
	newKey := res.Req1.Data.MusicKey
	if newKey == "" {
		log.Error("", "musicKey refresh failed, newKey is nil")
		return ""
	}
	log.Info("", "musicKey refresh success: "+uin+" newKey: "+newKey)
	return newKey
}
