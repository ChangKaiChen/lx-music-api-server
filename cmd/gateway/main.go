package main

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/config"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/middleware"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/router"
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/global"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	conf "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/cache"
	"github.com/hertz-contrib/cache/persist"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"strconv"
	"time"
)

const serviceName = consts.GateWayServiceName

var providers []provider.OtelProvider

func init() {
	Log := config.GetConf().Log
	logger.Init(serviceName, Log.Level, Log.Filepath)
	config.Init()
	providers = append(providers, rpc.Init()...)
}
func main() {
	log := logger.GetLogger()
	ctx := context.Background()
	defer func() {
		for _, p := range providers {
			err := p.Shutdown(ctx)
			if err != nil {
				log.Errorf("", "Failed to shutdown OpenTelemetry provider: %v", err)
			}
		}
	}()
	memoryStore := persist.NewMemoryStore(1000 * time.Second)
	middlewares := []app.HandlerFunc{
		middleware.RecoveryMW(),
		cache.NewCacheByRequestURI(memoryStore, 960*time.Second),
		middleware.SentinelMW(),
		middleware.LimiterMW(),
	}
	listenPort, err := config.GetPort()
	if err != nil {
		panic(err)
		return
	}
	opts := []conf.Option{
		server.WithHostPorts(":" + strconv.Itoa(listenPort)),
	}
	oTel := global.GetConf().OTel
	if oTel.Enable && oTel.Endpoint != "" {
		p := provider.NewOpenTelemetryProvider(
			provider.WithServiceName(serviceName),
			provider.WithInsecure(),
			provider.WithExportEndpoint(oTel.Endpoint),
			provider.WithHeaders(map[string]string{
				"X-ByteAPM-AppKey": oTel.Headers,
			}),
		)
		defer func(p provider.OtelProvider, ctx context.Context) {
			err = p.Shutdown(ctx)
			if err != nil {
				log.Errorf("", "Failed to shutdown OpenTelemetry provider: %v", err)
			}
		}(p, ctx)
		tracer, cfg := tracing.NewServerTracer()
		opts = append(opts, tracer)
		middlewares = append(middlewares, tracing.ServerMiddleware(cfg))
	}
	h := server.Default(opts...)
	h.Use(middlewares...)
	router.Register(h)
	h.Spin()
}
