package rpc

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/script"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ScriptHandler struct{}

func (s *ScriptHandler) Script(ctx context.Context, req *script.ScriptRequest) (*script.ScriptResponse, error) {
	log := logger.GetLogger()
	resp := &script.ScriptResponse{}
	checkUpdate := req.CheckUpdate
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get current file path")
	}
	sourceDir := filepath.Dir(filename)
	scriptPath := filepath.Join(sourceDir, "..", "lx-music-source.js")
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		resp.Success = false
		resp.Res.Code = consts.ServerErrorCode
		resp.Res.Msg = "Error reading script file"
		log.Errorf("", "Error reading script file: %v", err)
		return resp, err
	}
	hash := md5.Sum(content)
	hashString := hex.EncodeToString(hash[:])
	if checkUpdate != "" {
		if checkUpdate == hashString {
			resp.Success = true
			resp.Res = &script.Res{
				Code: consts.SuccessCode,
				Msg:  "This is already the latest version",
			}
			return resp, nil
		}
		resp.Success = true
		resp.Res = &script.Res{
			Code: consts.SuccessCode,
			Msg:  "This is not the latest version",
			Data: map[string]string{
				"updateMsg": "There is a new version of the lxmusic source",
				"updateUrl": fmt.Sprintf("%s/script?key=%s", req.Url, req.Key),
			},
		}
		return resp, nil
	}
	newContent := strings.Replace(
		string(content),
		"const API_KEY = ''",
		fmt.Sprintf("const API_KEY = '%s'", req.Key),
		1,
	)
	newContent = strings.Replace(newContent, "const API_URL = ''", fmt.Sprintf("const API_URL = '%s'", req.Url), 1)
	newContent = strings.Replace(newContent, "const SCRIPT_MD5 = ''", fmt.Sprintf("const SCRIPT_MD5 = '%s'", hashString), 1)
	resp.Success = true
	resp.Res = &script.Res{
		Content: newContent,
	}
	return resp, nil
}
