package selfout

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/internal/logic/selfout"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func SearchSelfoutHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.SearchSelfoutReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Fail(c, err)
			return
		}

		l := selfout.NewSearchSelfoutLogic(c.Request.Context(), svcCtx)
		resp, err := l.SearchSelfout(&req)
		if err != nil {
			httpresp.Fail(c, err)
		} else {
			//httpresp.JSONOk(c, resp)
			httpresp.OKJSON(c, resp)
		}
	}
}
