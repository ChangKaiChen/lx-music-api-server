package main

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/app/mg/config"
	"github.com/ChangKaiChen/lx-music-api-server/app/mg/rpc"
	"github.com/ChangKaiChen/lx-music-api-server/global"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/mg/mgservice"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/cache"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/utils"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"net"
	"runtime"
	"strconv"
)

const serviceName = consts.MgServiceName

func init() {
	Log := config.GetConf().Log
	logger.Init(serviceName, Log.Level, Log.Filepath)
	config.Init()
	global.Init()
	cache.Init()
}
func main() {
	var r registry.Registry
	var err error
	Etcd := global.GetConf().Etcd
	if Etcd.AuthEnable {
		r, err = etcd.NewEtcdRegistry([]string{Etcd.Addr}, etcd.WithAuthOpt(Etcd.Username, Etcd.Password))
	} else {
		r, err = etcd.NewEtcdRegistry([]string{Etcd.Addr})
	}
	if err != nil {
		panic(err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	address := ":" + strconv.Itoa(listenAddr)
	if global.GetConf().IsLocal {
		address = "localhost" + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts := []server.Option{
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithMaxConnIdleTime(Etcd.MaxIdleTimeout),
	}
	if runtime.GOOS != "windows" {
		opts = append(opts, server.WithMuxTransport())
	}
	oTel := global.GetConf().OTel
	if oTel.Enable && oTel.Endpoint != "" {
		ctx := context.Background()
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
			}
		}(p, ctx)
		opts = append(opts, server.WithSuite(tracing.NewServerSuite()))
	}
	svr := mgservice.NewServer(new(rpc.MgHandler), opts...)
	if err = svr.Run(); err != nil {
		panic(err)
	}
}
