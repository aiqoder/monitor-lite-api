package subscriber

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic/subscriber"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

func SubscriberListHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := subscriber.NewSubscriberListLogic(c.Request.Context(), svcCtx)
		resp, err := l.SubscriberList()
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OKJSON(c, resp)
		}
	}
}
