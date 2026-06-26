package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func BatchUpdateHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.BatchTvUpdateReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := logic.NewBatchUpdateLogic(c.Request.Context(), svcCtx)
		err := l.BatchUpdate(&req)
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OK(c)
		}
	}
}
