package proxy

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic/proxy"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func PtbtvHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.PtbTvReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := proxy.NewPtbtvLogic(c.Request.Context(), svcCtx)
		resp, err := l.Ptbtv(&req)
		if err != nil {
			httpresp.Error(c, err)
		} else {
			c.Redirect(http.StatusFound, resp)
		}
	}
}
