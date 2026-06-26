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

func IQiLuHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.IQiLuReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := proxy.NewIQiLuLogic(c.Request.Context(), svcCtx)
		resp, err := l.IQiLu(&req)
		if err != nil {
			httpresp.Error(c, err)
		} else {
			c.Redirect(http.StatusFound, resp)
		}
	}
}
