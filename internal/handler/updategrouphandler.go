package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

func UpdateGroupHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logic.NewUpdateGroupLogic(c.Request.Context(), svcCtx)
		err := l.UpdateGroup()
		if err != nil {
			httpresp.Error(c, err)
		} else {
			httpresp.OK(c)
		}
	}
}
