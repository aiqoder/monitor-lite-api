package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/pkg/common/mime"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func CustomOutHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CustomOutReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := logic.NewCustomOutLogic(c.Request.Context(), svcCtx)
		resp, err := l.CustomOut(&req)
		if err != nil {
			httpresp.Error(c, err)
		} else {
			c.Header("content-type", mime.TypeByExtension("txt")+"; charset=utf-8")
			_, _ = c.Writer.Write([]byte(resp))
		}
	}
}
