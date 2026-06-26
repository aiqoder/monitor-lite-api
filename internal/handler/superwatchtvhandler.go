package handler

import (
	httpresp "github.com/aiqoder/monitor-lite-api/internal/pkg/resp"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/bind"
	"github.com/gin-gonic/gin"
	"github.com/duke-git/lancet/v2/netutil"
	"slices"
	"github.com/aiqoder/monitor-lite-api/pkg/common/ipregion"
	"github.com/aiqoder/monitor-lite-api/pkg/common/mime"
	"github.com/aiqoder/monitor-lite-api/internal/logic"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

func SuperWatchTvHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req types.SuperWatchTvReq

		if err := bind.Parse(c, &req); err != nil {
			httpresp.Error(c, err)
			return
		}

		l := logic.NewSuperWatchTvLogic(c.Request.Context(), svcCtx)
		ip := netutil.GetRequestPublicIp(c.Request)
		region, err2 := ipregion.Ip2regionSearch(ip)

		if err2 != nil {
			httpresp.Error(c, err2)
			return
		}

		if region.Country == "中国" && !slices.Contains([]string{"台湾省", "香港", "澳门"}, region.Province) {
			resp, err := l.SuperWatchTv(&req)
			if err != nil {
				httpresp.Error(c, err)
			} else {
				//httpresp.OKJSON(c, resp)
				c.Header("content-type", mime.TypeByExtension("txt")+"; charset=utf-8")
				//c.Header("Cache-Control", "public,max-age=30")
				_, _ = c.Writer.Write([]byte(resp))
			}
		} else {
			resp, err := l.SuperOutWatchTv(&req)
			if err != nil {
				httpresp.Error(c, err)
			} else {
				//httpresp.OKJSON(c, resp)
				c.Header("content-type", mime.TypeByExtension("txt")+"; charset=utf-8")
				//c.Header("Cache-Control", "public,max-age=30")
				_, _ = c.Writer.Write([]byte(resp))
			}
		}
	}
}
