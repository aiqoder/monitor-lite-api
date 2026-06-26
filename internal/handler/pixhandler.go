package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

func PixHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logic.NewPixLogic(c.Request.Context(), svcCtx)
		resp, err := l.Pix()
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OKJSON(c, resp)
		}
	}
}
