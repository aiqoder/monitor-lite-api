package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func SyncTVDataHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.SyncTVDataReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := logic.NewSyncTVDataLogic(c.Request.Context(), svcCtx)
		err := l.SyncTVData(&req)
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OK(c)
		}
	}
}
