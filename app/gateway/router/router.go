package router

import (
	apikg "github.com/ChangKaiChen/lx-music-api-server/app/gateway/router/api/kg"
	apikw "github.com/ChangKaiChen/lx-music-api-server/app/gateway/router/api/kw"
	apimg "github.com/ChangKaiChen/lx-music-api-server/app/gateway/router/api/mg"
	apiscript "github.com/ChangKaiChen/lx-music-api-server/app/gateway/router/api/script"
	apitx "github.com/ChangKaiChen/lx-music-api-server/app/gateway/router/api/tx"
	apiwy "github.com/ChangKaiChen/lx-music-api-server/app/gateway/router/api/wy"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register(h *server.Hertz) {
	apiscript.Register(h)
	urlGroup := h.Group("/url")
	apiwy.Register(urlGroup.Group("/wy"))
	apikg.Register(urlGroup.Group("/kg"))
	apitx.Register(urlGroup.Group("/tx"))
	apimg.Register(urlGroup.Group("/mg"))
	apikw.Register(urlGroup.Group("/kw"))
}
