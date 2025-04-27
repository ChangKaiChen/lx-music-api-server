package refresh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ChangKaiChen/lx-music-api-server/app/wy/crypto"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Header struct {
	OSVer       string `json:"osver"`
	DeviceID    string `json:"deviceId"`
	OS          string `json:"os"`
	AppVer      string `json:"appver"`
	VersionCode string `json:"versioncode"`
	MobileName  string `json:"mobilename"`
	BuildVer    string `json:"buildver"`
	Resolution  string `json:"resolution"`
	Csrf        string `json:"__csrf"`
	Channel     string `json:"channel"`
	RequestID   string `json:"requestId"`
	MusicU      string `json:"MUSIC_U"`
	MusicA      string `json:"MUSIC_A"`
}
type Res struct {
	BizCode string `json:"bizCode"`
	Code    int    `json:"code"`
}

func CookieRefresh(cookie string) bool {
	log := logger.GetLogger()
	co := cookieStr2Dict(cookie)
	var header Header
	header.OSVer = "17.4.1"
	if co["osver"] != "" {
		header.OSVer = co["osver"]
	}
	if co["deviceId"] != "" {
		header.DeviceID = co["deviceId"]
	}
	header.OS = "ios"
	if co["os"] != "" {
		header.OS = co["os"]
	}
	if header.OS != "pc" {
		header.AppVer = "9.0.65"
	}
	if co["appver"] != "" {
		header.AppVer = co["appver"]
	}
	header.VersionCode = "140"
	if co["versioncode"] != "" {
		header.VersionCode = co["versioncode"]
	}
	if co["mobilename"] != "" {
		header.MobileName = co["mobilename"]
	}
	header.BuildVer = strconv.FormatInt(time.Now().Unix(), 10)[:10]
	if co["buildver"] != "" {
		header.BuildVer = co["buildver"]
	}
	header.Resolution = "1920x1080"
	if co["resolution"] != "" {
		header.Resolution = co["resolution"]
	}
	if co["__csrf"] != "" {
		header.Csrf = co["__csrf"]
	}
	if co["channel"] != "" {
		header.Channel = co["channel"]
	}
	header.RequestID = strconv.FormatInt(time.Now().UnixNano()/1000000, 10)[:13] + "_" + fmt.Sprintf("%04d", rand.Intn(10000))
	if co["MUSIC_U"] != "" {
		header.MusicU = co["MUSIC_U"]
	}
	if co["MUSIC_A"] != "" {
		header.MusicA = co["MUSIC_A"]
	}
	jsonBody, err := json.Marshal(map[string]any{"header": header, "e_r": false})
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return false
	}
	jsonStr := string(jsonBody)
	baseUrl := "http://interface.music.163.com/eapi/"
	path := "/api/login/token/refresh"
	formMap := crypto.EApiEncrypt(path, jsonStr)
	form := url.Values{}
	form.Add("params", formMap["params"])
	r, err := http.NewRequest("POST", baseUrl+path[5:], bytes.NewBufferString(form.Encode()))
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return false
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Cookie", cookieDict2Str(header))
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(r)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return false
	}
	defer func() {
		err = httpResp.Body.Close()
		if err != nil {
		}
	}()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return false
	}
	var res Res
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Errorf("", "refresh error: %v", err.Error())
		return false
	}
	if res.Code != 200 {
		log.Error("", "refresh failed")
		return false
	}
	log.Info("", "cookie refresh success")
	return true
}
func cookieStr2Dict(cookieStr string) map[string]string {
	cookieDict := make(map[string]string)
	parts := strings.Split(cookieStr, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			cookieDict[kv[0]] = kv[1]
		}
	}
	return cookieDict
}
func cookieDict2Str(header Header) string {
	var cookieStr string
	if header.OSVer != "" {
		cookieStr += fmt.Sprintf("osver=%s; ", header.OSVer)
	}
	if header.DeviceID != "" {
		cookieStr += fmt.Sprintf("deviceId=%s; ", header.DeviceID)
	}
	if header.OS != "" {
		cookieStr += fmt.Sprintf("os=%s; ", header.OS)
	}
	if header.AppVer != "" {
		cookieStr += fmt.Sprintf("appver=%s; ", header.AppVer)
	}
	if header.VersionCode != "" {
		cookieStr += fmt.Sprintf("versioncode=%s; ", header.VersionCode)
	}
	if header.MobileName != "" {
		cookieStr += fmt.Sprintf("mobilename=%s; ", header.MobileName)
	}
	if header.BuildVer != "" {
		cookieStr += fmt.Sprintf("buildver=%s; ", header.BuildVer)
	}
	if header.Resolution != "" {
		cookieStr += fmt.Sprintf("resolution=%s; ", header.Resolution)
	}
	if header.Csrf != "" {
		cookieStr += fmt.Sprintf("__csrf=%s; ", header.Csrf)
	}
	if header.Channel != "" {
		cookieStr += fmt.Sprintf("channel=%s; ", header.Channel)
	}
	if header.MusicU != "" {
		cookieStr += fmt.Sprintf("MUSIC_U=%s; ", header.MusicU)
	}
	if header.MusicA != "" {
		cookieStr += fmt.Sprintf("MUSIC_A=%s; ", header.MusicA)
	}
	return cookieStr
}
