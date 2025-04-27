package mg

import (
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/handler/mg"
	"github.com/cloudwego/hertz/pkg/route"
)

func Register(group *route.RouterGroup) {
	group.GET("/:uid/:quality", handler.MgMusicUrl)
}
