package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/duke-git/lancet/v2/netutil"
	"slices"
	"github.com/aiqoder/monitor-lite-api/pkg/common/ipregion"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func SearchHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.TvJsonReq
		if err := bind.Parse(c, &req); err != nil {
			httpresp.Fail(c, err)
			return
		}

		l := logic.NewSearchLogic(c.Request.Context(), svcCtx)

		ip := netutil.GetRequestPublicIp(c.Request)
		region, err2 := ipregion.Ip2regionSearch(ip)

		if err2 != nil {
			httpresp.Error(c, err2)
			return
		}

		resp, err := l.Search(&req, region.Country == "中国" && !slices.Contains([]string{"台湾省", "香港", "澳门"}, region.Province))
		if err != nil {
			httpresp.Fail(c, err)
		} else {
			httpresp.JSONOk(c, resp)
		}
	}
}
