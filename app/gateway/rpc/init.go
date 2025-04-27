package rpc

import (
	"github.com/ChangKaiChen/lx-music-api-server/global"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/auth/authservice"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/kg/kgservice"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/kw/kwservice"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/mg/mgservice"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/script/scriptservice"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/tx/txservice"
	"github.com/ChangKaiChen/lx-music-api-server/kitex_gen/wy/wyservice"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"runtime"
)

var (
	AuthClient   *authservice.Client
	ScriptClient *scriptservice.Client
	TxClient     *txservice.Client
	WyClient     *wyservice.Client
	KgClient     *kgservice.Client
	KwClient     *kwservice.Client
	MgClient     *mgservice.Client
)
var (
	MaxIdleTimeout    = global.GetConf().Etcd.MaxIdleTimeout
	MinIdlePerAddress = global.GetConf().Etcd.MinIdlePerAddress
	opts              []client.Option
	providers         []provider.OtelProvider
)

func Init() []provider.OtelProvider {
	Etcd := global.GetConf().Etcd
	var r discovery.Resolver
	var err error
	if Etcd.AuthEnable {
		r, err = etcd.NewEtcdResolver([]string{Etcd.Addr}, etcd.WithAuthOpt(Etcd.Username, Etcd.Password))
	} else {
		r, err = etcd.NewEtcdResolver([]string{Etcd.Addr})
	}
	if err != nil {
		panic(err)
	}
	opts = []client.Option{
		client.WithResolver(r), client.WithTransportProtocol(transport.TTHeader), client.WithSuite(tracing.NewClientSuite()),
	}
	if runtime.GOOS == "windows" {
		opts = append(opts, client.WithLongConnection(connpool.IdleConfig{MaxIdleTimeout: MaxIdleTimeout, MinIdlePerAddress: MinIdlePerAddress}))
	} else {
		opts = append(opts, client.WithMuxConnection(consts.MuxConnectionNum))
	}
	AuthClient = AuthRPCInit()
	ScriptClient = ScriptRPCInit()
	TxClient = TxRPCInit()
	WyClient = WyRPCInit()
	KgClient = KgRPCInit()
	KwClient = KwRPCInit()
	MgClient = MgRPCInit()
	return providers
}
func AuthRPCInit() *authservice.Client {
	oTelProvider(consts.AuthServiceName)
	authClient, err := authservice.NewClient(consts.AuthServiceName, opts...)
	if err != nil {
		panic(err)
	}
	return &authClient
}
func WyRPCInit() *wyservice.Client {
	oTelProvider(consts.WyServiceName)
	clientWy, err := wyservice.NewClient(consts.WyServiceName, opts...)
	if err != nil {
		panic(err)
	}
	return &clientWy
}
func ScriptRPCInit() *scriptservice.Client {
	oTelProvider(consts.ScriptServiceName)
	clientScript, err := scriptservice.NewClient(consts.ScriptServiceName, opts...)
	if err != nil {
		panic(err)
	}
	return &clientScript
}
func TxRPCInit() *txservice.Client {
	oTelProvider(consts.TxServiceName)
	clientTx, err := txservice.NewClient(consts.TxServiceName, opts...)
	if err != nil {
		panic(err)
	}
	return &clientTx
}
func KgRPCInit() *kgservice.Client {
	oTelProvider(consts.KgServiceName)
	clientKg, err := kgservice.NewClient(consts.KgServiceName, opts...)
	if err != nil {
		panic(err)
	}
	return &clientKg
}
func KwRPCInit() *kwservice.Client {
	oTelProvider(consts.KwServiceName)
	clientKw, err := kwservice.NewClient(consts.KwServiceName, opts...)
	if err != nil {
		panic(err)
	}
	return &clientKw
}
func MgRPCInit() *mgservice.Client {
	oTelProvider(consts.MgServiceName)
	clientMg, err := mgservice.NewClient(consts.MgServiceName, opts...)
	if err != nil {
		panic(err)
	}
	return &clientMg
}
func oTelProvider(serviceName string) {
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
		providers = append(providers, p)
	}
}
