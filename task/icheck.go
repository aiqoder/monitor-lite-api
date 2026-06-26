package task

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"strconv"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/ffmpeg"
	"github.com/aiqoder/monitor-lite-api/pkg/common/tools"
	"github.com/aiqoder/monitor-lite-api/atools"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/model"
)

// ICheck 用于定时检测数据库当中的直播源
func ICheck(ctx *svc.ServiceContext) func() {
	return func() {
		if ctx.SettingModel.Value("autoCheck") == "0" {
			log.Error("未开启自动测源")
			return
		}

		offset, err := ctx.TvModel.NextTvOffset()
		if err != nil {
			log.Error(err)
		}

		for i := 0; i < len(offset); i++ {
			msg := fmt.Sprintf("正在自动检测，%d,%s", i+1, offset[i].Name)
			atools.WsMsg(ctx, "后台自动检测", msg)
			time.Sleep(200 * time.Millisecond)

			tv := offset[i]

			// 删除日本字节目
			if tools.IsJapanString(tv.Name) {
				ctx.TvModel.Delete([]string{strconv.FormatUint(tv.ID, 10)})
				atools.WsMsg(ctx, "健康助手提示", fmt.Sprintf("疑似发现日本字：%s，节目已进入垃圾桶", tv.Name))
				continue
			}

			// 删除疑似H节目 过滤无意义名称
			continueStr := []string{"三级", "强奸", "车震"}
			if tools.StrInContains(tv.Name, continueStr) {
				ctx.TvModel.Delete([]string{strconv.FormatUint(tv.ID, 10)})
				atools.WsMsg(ctx, "健康助手提示", fmt.Sprintf("疑似触发非法关键词：%s，节目已进入垃圾桶", tv.Name))
				continue
			}

			// 悄悄删除无意义的节目
			continueStr2 := []string{"直播室", "强奸", "车震"}
			if tools.StrInContains(tv.Name, continueStr2) {
				ctx.TvModel.Delete([]string{strconv.FormatUint(tv.ID, 10)})
				continue
			}

			info, err := ffmpeg.GetOnlineVideoInfo(tv.Url)
			if err != nil {
				// 删除失败的源
				autoDelCount := ctx.SettingModel.Value("autoDelCount")
				to, _ := convertor.ToInt(autoDelCount)
				if (tv.FailCount + 1) >= to {
					ctx.TvModel.Delete([]string{strconv.FormatUint(tv.ID, 10)})
					continue
				}
				_ = ctx.TvModel.Update(model.Tv{ID: tv.ID}, true)
				continue
			}
			_ = ctx.TvModel.Update(model.Tv{ID: tv.ID, Speed: info.Speed, Width: info.Width, Height: info.Height}, false)
		}
	}
}
