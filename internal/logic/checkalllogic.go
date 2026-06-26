package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"strconv"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/ffmpeg"
	"github.com/aiqoder/monitor-lite-api/atools"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
	"github.com/aiqoder/monitor-lite-api/model"
)

type CheckAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckAllLogic {
	return &CheckAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var updatedCheckFlag = false

func (l *CheckAllLogic) autoCheck(tv model.Tv) {
	info, err := ffmpeg.GetOnlineVideoInfo(tv.Url)
	if err != nil {
		// 删除失败的源
		autoDelCount := l.svcCtx.SettingModel.Value("autoDelCount")
		to, _ := convertor.ToInt(autoDelCount)
		if (tv.FailCount + 1) >= to {
			l.svcCtx.TvModel.Delete([]string{strconv.FormatUint(tv.ID, 10)})
			return
		}
		// 更新源
		_ = l.svcCtx.TvModel.Update(model.Tv{ID: tv.ID}, true)
		return
	}
	_ = l.svcCtx.TvModel.Update(model.Tv{ID: tv.ID, Speed: info.Speed, Width: info.Width, Height: info.Height}, false)
}

func (l *CheckAllLogic) CheckAll(req *types.CheckAllReq) error {
	if updatedCheckFlag {
		return errors.New("后台正在检测中，请勿重复下发指令！")
	}
	updatedCheckFlag = true

	go func() {
		offset, err := l.svcCtx.TvModel.TableCheckList(req.Type, req.Extra)
		if err != nil {
			atools.TipsQueen.PushBack(err.Error())
			return
		}

		length := len(offset)
		for i := 0; i < length; i++ {
			time.Sleep(50 * time.Millisecond)
			tv := offset[i]
			msg := fmt.Sprintf("正在检测：(%s),当前进度（%d/%d）", tv.Name, i+1, length)
			atools.WsMsg(l.svcCtx, "一键监测", msg)
			l.autoCheck(tv)
		}
	}()

	updatedCheckFlag = false
	return nil
}
