package logic

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"strings"
	"github.com/aiqoder/monitor-lite-api/pkg/common/tools"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
	"github.com/aiqoder/monitor-lite-api/model"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/rulecache"
)

type CustomOutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCustomOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CustomOutLogic {
	return &CustomOutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CustomOutLogic) CustomOut(req *types.CustomOutReq) (resp string, err error) {
	self, err := l.svcCtx.SelfoutModel.FirstOneByKey(req.Path)
	if err != nil {
		return "", err
	}

	//if fileutil.IsExist("/.dockerenv") {
	//	authCode := l.svcCtx.SettingModel.Value("authCode")
	//}
	cacheGroupText, ok := l.svcCtx.Cache.Get("cache-group")
	unknownGroup := l.svcCtx.SettingModel.Value("unknownGroup") == "1"

	if ok {
		str := cacheGroupText.(string)
		return str, nil
	}

	toString := l.svcCtx.SettingModel.RuleContent()

	groups := getGroup(toString)

	resultGroup := ""

	// 查询黑名单
	blackList := strings.Split(l.svcCtx.SettingModel.Value("blackList"), "#")
	// 查询数据
	list, _ := l.svcCtx.TvModel.TableList()

	var list2 []model.Tv

	for _, v := range list {
		if v.Speed < uint64(self.Speed) {
			continue
		}

		if v.Width < uint64(self.Width) {
			continue
		}

		if v.Height < uint64(self.Height) {
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

	return resultGroup, nil
}
