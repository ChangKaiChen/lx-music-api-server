package rpc

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/auth"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func AuthRPC(ctx context.Context, c *app.RequestContext, key string) (*auth.AuthResponse, error) {
	log := logger.GetLogger()
	resp := &auth.AuthResponse{}
	reqAuth := &auth.AuthRequest{
		AuthKey: key,
	}
	resp, err := (*rpc.AuthClient).Auth(ctx, reqAuth)
	if err != nil {
		response.ServerErrorResponse(c, "Internal server error")
		log.Errorf("", "AuthRPC call failed: %v", err)
		return resp, err
	}
	return resp, nil
}
