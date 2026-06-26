package epg

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic/epg"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

func EpgCollectHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := epg.NewEpgCollectLogic(c.Request.Context(), svcCtx)
		err := l.EpgCollect()
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OK(c)
		}
	}
}
