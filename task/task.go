package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/robfig/cron/v3"
	"github.com/aiqoder/monitor-lite-api/atools"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/aigroup"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/model"
	"github.com/aiqoder/monitor-lite-api/utils"
)

var isUpdateGroup = false

func UpdateGroup(svcCtx *svc.ServiceContext, mode string) func() {
	return func() {
		if isUpdateGroup {
			atools.WsMsg(svcCtx, "后台更新分组", "后台正在更新分组，请不要着急！！！")
			return
		}
		isUpdateGroup = true
		defer func() { isUpdateGroup = false }()

		if mode == "auto" && !aigroup.NewClient(utils.AISettingsFromCtx(svcCtx)).Enabled() {
			return
		}

		start := time.Now().UnixMilli()
		tvList, _ := svcCtx.TvModel.TableListGroupName()

		var pending []model.Tv
		for i := range tvList {
			tv := tvList[i]
			if len(strings.ReplaceAll(tv.Group, " ", "")) > 0 {
				continue
			}
			if validator.IsNumberStr(tv.Name) {
				continue
			}
			pending = append(pending, tv)
		}

		if len(pending) == 0 {
			return
		}

		err := utils.UpdateGroupByAI(svcCtx, pending)
		if err != nil {
			log.Error("AI 分组失败:", err)
			atools.WsMsg(svcCtx, "后台更新分组", err.Error())
		} else {
			svcCtx.UpdatePrompt()
			atools.WsMsg(svcCtx, "后台更新分组", fmt.Sprintf("AI 分组完成，处理 %d 条，耗时 %d ms", len(pending), time.Now().UnixMilli()-start))
		}
	}
}

func Start(ctx *svc.ServiceContext) {
	search := ISearch(ctx)
	check := ICheck(ctx)
	grab := Igrab(ctx)
	update := UpdateGroup(ctx, "auto")
	epg := IEpg(ctx)

	go func() {
		grab()
	}()

	go func() {
		update()
	}()

	cronTab := cron.New(cron.WithSeconds())
	_, _ = cronTab.AddFunc("@every 120m", update)
	_, _ = cronTab.AddFunc("@every 10m", check)
	_, _ = cronTab.AddFunc("@every 15m", search)
	_, _ = cronTab.AddFunc("0 8 * * *", epg)
	cronTab.Start()
}
