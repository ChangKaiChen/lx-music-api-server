package wy

import (
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/handler/wy"
	"github.com/cloudwego/hertz/pkg/route"
)

func Register(group *route.RouterGroup) {
	group.GET("/:uid/:quality", handler.WyMusicUrl)
}
