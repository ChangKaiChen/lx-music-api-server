package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/auth/config"
	"github.com/ChangKaiChen/lx-music-api-server/app/auth/crypto"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/auth"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
)

type AuthHandler struct{}

func (s *AuthHandler) Auth(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	log := logger.GetLogger()
	resp := &auth.AuthResponse{}
	key := req.AuthKey
	decryptData, err := crypto.Decrypt(key)
	if err != nil {
		log.Errorf("", "Key: %v, Decrypt error: %v", key, err)
		return resp, err
	}
	token := decryptData
	tokens := config.GetConf().WhiteList
	for _, item := range tokens {
		if item.Token == token {
			resp.Success = true
			resp.Res = map[string]string{
				"token": token,
			}
			log.Log("", key, "", "", "", "Success")
			return resp, nil
		}
	}
	resp.Res = map[string]string{
		"msg": "Unauthorized",
	}
	log.Log("", key, "", "", "", "Unauthorized")
	return resp, nil
}
