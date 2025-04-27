package kw

import (
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/handler/kw"
	"github.com/cloudwego/hertz/pkg/route"
)

func Register(group *route.RouterGroup) {
	group.GET("/:uid/:quality", handler.KwMusicUrl)
}
