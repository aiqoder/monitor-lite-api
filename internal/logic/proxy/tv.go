package proxy

import (
	"github.com/aiqoder/monitor-lite-api/model"
)

var ProxyTvs = []model.Tv{
	{
		Name: "凤凰资讯",
		Url:  "/v1/proxy/fengshows?t=zx",
	},
	{
		Name: "凤凰中文",
		Url:  "/v1/proxy/fengshows?t=zw",
	},
	{
		Name: "香港台",
		Url:  "/v1/proxy/fengshows?t=xg",
	},
	{
		Name: "莆田1套",
		Url:  "/v1/proxy/ptbtv?t=pt1",
	},
	{
		Name: "莆田2套",
		Url:  "/v1/proxy/ptbtv?t=pt2",
	},
	{
		Name: "仙游电视",
		Url:  "/v1/proxy/ptbtv?t=xy",
	},
	// 爱齐鲁
	{
		Name: "山东齐鲁频道",
		Url:  "/v1/proxy/iqilu?t=ql",
	},
	{
		Name: "山东体育休闲频道",
		Url:  "/v1/proxy/iqilu?t=ty",
	},
	{
		Name: "山东生活频道",
		Url:  "/v1/proxy/iqilu?t=sh",
	},
	{
		Name: "山东综艺频道",
		Url:  "/v1/proxy/iqilu?t=zy",
	},
	{
		Name: "山东农科频道",
		Url:  "/v1/proxy/iqilu?t=nk",
	},
	{
		Name: "山东文旅频道",
		Url:  "/v1/proxy/iqilu?t=wl",
	},
}
