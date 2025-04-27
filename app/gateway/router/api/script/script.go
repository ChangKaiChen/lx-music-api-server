package script

import (
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/handler/script"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register(h *server.Hertz) {
	h.GET("/script", handler.Script)
}
