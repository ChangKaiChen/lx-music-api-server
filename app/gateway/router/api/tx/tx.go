package tx

import (
	"github.com/ChangKaiChen/lx-music-api-server/app/gateway/handler/tx"
	"github.com/cloudwego/hertz/pkg/route"
)

func Register(group *route.RouterGroup) {
	group.GET("/:uid/:quality", handler.TxMusicUrl)
}
