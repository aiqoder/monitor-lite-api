package subscriber

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic/subscriber"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func SubscriberUpdateHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.Subscriber
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := subscriber.NewSubscriberUpdateLogic(c.Request.Context(), svcCtx)
		resp, err := l.SubscriberUpdate(&req)
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OKJSON(c, resp)
		}
	}
}
