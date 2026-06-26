package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func ChangePasswordHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ChangePasswordReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := logic.NewChangePasswordLogic(c.Request.Context(), svcCtx)
		if err := l.ChangePassword(&req); err != nil {
			httpresp.Error(c, err)
			return
		}
		httpresp.OKJSON(c, &types.ChangePasswordResp{})
	}
}
