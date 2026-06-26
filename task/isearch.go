package task

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/aiqoder/monitor-lite-api/pkg/common/ffmpeg"
	"github.com/aiqoder/monitor-lite-api/atools"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/prompt"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/model"
	"github.com/aiqoder/monitor-lite-api/utils"
)

func ISearch(ctx *svc.ServiceContext) func() {
	q := list.New()
	return func() {
		if ctx.SettingModel.Value("autoSearch") == "0" {
			log.Error("未开启自动找源")
			return
		}
		if q.Len() == 0 {
			cfg, err := prompt.Parse(ctx.SettingModel.RuleContent())
			if err != nil {
				log.Error(err)
				return
			}
			seen := make(map[string]struct{})
			push := func(name string) {
				name = strings.TrimSpace(name)
				if name == "" {
					return
				}
				if _, ok := seen[name]; ok {
					return
				}
				seen[name] = struct{}{}
				q.PushBack(name)
			}
			for _, name := range cfg.AllChannelNames() {
				push(name)
			}
		}

		find := func(key any) (int, bool) {
			tvName := key.(string)
			pureName := strings.ReplaceAll(tvName, " ", "")
			pureName = strings.ReplaceAll(pureName, "\r", "")

			if len(pureName) == 0 {
				return 0, false
			}
			log.Info("开始新一轮搜索，触发关键字：", pureName)
			msg := fmt.Sprintf("正在触发关键字，%s", pureName)
			atools.WsMsg(ctx, "后台搜索检测", msg)

			tvs, err := utils.Search(pureName, "name")
			if err != nil {
				log.Error(err)
				return 0, false
			}

			if len(tvs) > 100 {
				tvs = tvs[:100]
			}

			length := len(tvs)
			if length == 0 {
				return 0, false
			}

			log.Info("找到结果数量：", length)
			msg = fmt.Sprintf("正在触发关键字[%s]，搜索到结果数量%d", pureName, length)
			atools.WsMsg(ctx, "后台搜索检测", msg)
			for i := 0; i < length; i++ {
				tv := tvs[i]
				info, err := ffmpeg.GetOnlineVideoInfo(tv.Url)
				if err != nil {
					continue
				}
				_ = ctx.TvModel.Save(model.Tv{Name: tv.Name, Url: tv.Url, Speed: info.Speed, Width: info.Width, Height: info.Height}, false)
			}

			return length, true
		}

		findFlag := false
		for i := 0; i < 50; i++ {
			if findFlag {
				break
			}
			if q.Len() == 0 {
				break
			}
			key := q.Remove(q.Front())
			num, _ := find(key)
			if num == 0 {
				continue
			}
			findFlag = true
		}
	}
}
