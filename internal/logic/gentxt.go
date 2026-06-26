package logic

import (
	"fmt"
	"os"
	"strings"
	"github.com/aiqoder/monitor-lite-api/pkg/common/tools"
	"github.com/aiqoder/monitor-lite-api/internal/logic/proxy"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/rulecache"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/model"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
)

// domestic false 表示国外， true 表示国内
func genTxt(svcCtx *svc.ServiceContext, domestic bool) string {
	//if fileutil.IsExist("/.dockerenv") {
	//	authCode := svcCtx.SettingModel.Value("authCode")
	//}
	cacheGroupText, ok := svcCtx.Cache.Get("cache-group")
	unknownGroup := svcCtx.SettingModel.Value("unknownGroup") == "1"

	if ok {
		str := cacheGroupText.(string)
		return str
	}

	toString := svcCtx.SettingModel.RuleContent()

	groups := getGroup(toString)

	resultGroup := ""

	// 查询黑名单
	blackList := strings.Split(svcCtx.SettingModel.Value("blackList"), "#")
	// 查询数据
	list, _ := svcCtx.TvModel.TableList()

	var list2 []model.Tv

	for _, v := range list {
		if v.Speed == 0 {
			continue
		}

		// 宽高小于0不进行组装
		if v.Width <= 0 {
			continue
		}

		// 失败的源不进行展示
		if v.FailCount > 0 {
			continue
		}

		// 屏蔽的源将不显示
		if tools.StrInContains(v.Url, blackList) {
			continue
		}

		gs := strings.Split(v.Group, ",")
		for i := 0; i < len(gs); i++ {
			group := gs[i]
			if strutil.IsBlank(group) {
				group = "未知分组"
			}
			list2 = append(list2, model.Tv{
				Name:        v.Name,
				DisplayName: v.DisplayName,
				Url:         v.Url,
				UpdateTime:  v.UpdateTime,
				Height:      v.Height,
				Width:       v.Width,
				Speed:       v.Speed,
				FailCount:   v.FailCount,
				Group:       group,
			})
		}
	}

	with := slice.GroupWith(list2, func(item model.Tv) string {
		return item.Group
	})

	for i := 0; i < len(groups); i++ {
		groupName := groups[i]

		// 国内不允许出现港澳台的分组
		if domestic {
			if strings.Contains(groupName, "香港") {
				continue
			}

			if strings.Contains(groupName, "台湾") {
				continue
			}

			if strings.Contains(groupName, "澳门") {
				continue
			}
		}

		// 将 group 当中对应的电视名称分开
		groupTvName, _ := rulecache.GetGroupChannels(groupName)
		groupTvNames := resolveGroupChannelNames(groupTvName, with[groupName])
		if len(groupTvNames) == 0 {
			continue
		}
		var groupTvs []string
		for i := 0; i < len(groupTvNames); i++ {
			rightTvName := groupTvNames[i]

			var urls []string

			tvs := with[groupName]
			length := len(tvs)

			for i := 0; i < length; i++ {
				tv := tvs[i]

				if tv.DisplayName == rightTvName {
					urls = append(urls, tv.Url)
				}
			}

			if len(urls) > 0 {
				groupTvs = append(groupTvs, fmt.Sprintf("%s,%s", rightTvName, strings.Join(urls, "#")))
			}
		}

		if len(groupTvs) > 0 {
			resultGroup += fmt.Sprintf("%s,#genre#\n%s\n\n", groupName, strings.Join(groupTvs, "\n"))
		}
	}

	if unknownGroup {
		tvs := with["未知分组"]
		tvWithByName := slice.GroupWith(tvs, func(item model.Tv) string {
			return item.Name
		})

		var groupTvs []string
		for tvName, tvs := range tvWithByName {
			urls := slice.Map(tvs, func(index int, item model.Tv) string {
				return item.Url
			})
			groupTvs = append(groupTvs, fmt.Sprintf("%s,%s", tvName, strings.Join(urls, "#")))
		}

		if len(groupTvs) > 0 {
			resultGroup += fmt.Sprintf("%s,#genre#\n%s\n\n", "未知分组", strings.Join(groupTvs, "\n"))
		}
	}

	svcCtx.Cache.Set("cache-group", resultGroup)
	return resultGroup
}

func genM3u(svcCtx *svc.ServiceContext) string {
	cacheGroupText, ok := svcCtx.Cache.Get("cache-group")
	unknownGroup := svcCtx.SettingModel.Value("unknownGroup") == "1"

	if ok {
		str := cacheGroupText.(string)
		return str
	}

	toString := svcCtx.SettingModel.RuleContent()
	groups := getGroup(toString)

	resultGroup := "#EXTM3U x-tvg-url=\"https://epg.112114.xyz/pp.xml.gz,https://assets.livednow.com/epg.xml\"\n"

	// 查询黑名单
	blackList := strings.Split(svcCtx.SettingModel.Value("blackList"), "#")
	// 查询数据
	list, _ := svcCtx.TvModel.TableList()

	var list2 []model.Tv

	for _, v := range list {
		if v.Speed == 0 {
			continue
		}

		// 宽高小于0不进行组装
		if v.Width <= 0 {
			continue
		}

		// 失败的源不进行展示
		if v.FailCount > 0 {
			continue
		}

		// 屏蔽的源将不显示
		if tools.StrInContains(v.Url, blackList) {
			continue
		}

		gs := strings.Split(v.Group, ",")
		for i := 0; i < len(gs); i++ {
			group := gs[i]
			if strutil.IsBlank(group) {
				group = "未知分组"
			}
			list2 = append(list2, model.Tv{
				Name:        v.Name,
				DisplayName: v.DisplayName,
				Url:         v.Url,
				UpdateTime:  v.UpdateTime,
				Height:      v.Height,
				Width:       v.Width,
				Speed:       v.Speed,
				FailCount:   v.FailCount,
				Group:       group,
			})
		}
	}

	if os.Getenv("PROXY") == "1" {
		for i := 0; i < len(proxy.ProxyTvs); i++ {
			tv := &proxy.ProxyTvs[i]
			tv.Url = os.Getenv("PROXY_HOST") + tv.Url
		}
		list2 = append(list2, proxy.ProxyTvs...)
		os.Setenv("PROXY_HOST", "")
	}

	with := slice.GroupWith(list2, func(item model.Tv) string {
		return item.Group
	})

	for i := 0; i < len(groups); i++ {
		groupName := groups[i]

		// 将 group 当中对应的电视名称分开
		groupTvName, _ := rulecache.GetGroupChannels(groupName)
		groupTvNames := resolveGroupChannelNames(groupTvName, with[groupName])
		if len(groupTvNames) == 0 {
			continue
		}
		var groupTvs []string
		for i := 0; i < len(groupTvNames); i++ {
			rightTvName := groupTvNames[i]

			var urls []string

			tvs := with[groupName]
			length := len(tvs)

			for i := 0; i < length; i++ {
				tv := tvs[i]

				if tv.DisplayName == rightTvName {
					urls = append(urls, tv.Url)
				}

				if len(urls) > 20 {
					break
				}
			}

			if len(urls) > 0 {
				for _, url := range urls {
					rightNameLine := fmt.Sprintf("#EXTINF:-1 tvg-name=\"%s\" tvg-logo=\"\" group-title=\"%s\",%s", rightTvName, groupName, rightTvName)
					groupTvs = append(groupTvs, fmt.Sprintf("%s\n%s\n", rightNameLine, url))
				}
				resultGroup += strings.Join(groupTvs, "\n")
			}
		}
	}

	if unknownGroup {
		tvs := with["未知分组"]
		tvWithByName := slice.GroupWith(tvs, func(item model.Tv) string {
			return item.Name
		})

		var groupTvs []string
		for tvName, tvs := range tvWithByName {
			urls := slice.Map(tvs, func(index int, item model.Tv) string {
				return item.Url
			})
			for _, url := range urls {
				rightNameLine := fmt.Sprintf("#EXTINF:-1 tvg-name=\"%s\" tvg-logo=\"\" group-title=\"%s\",%s", tvName, "未知分组", tvName)
				groupTvs = append(groupTvs, fmt.Sprintf("%s\n%s\n", rightNameLine, url))
			}
			resultGroup += strings.Join(groupTvs, "\n")
		}
	}

	svcCtx.Cache.Set("cache-group", resultGroup)
	return resultGroup
}
