package epg

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic/epg"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func EpgSearchHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.EpgSearchReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := epg.NewEpgSearchLogic(c.Request.Context(), svcCtx)
		resp, err := l.EpgSearch(&req)
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OKJSON(c, resp)
		}
	}
}
