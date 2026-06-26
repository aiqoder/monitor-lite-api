package logic

import (
	"context"
	"github.com/duke-git/lancet/v2/convertor"
	"strconv"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
	"github.com/aiqoder/monitor-lite-api/model"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.TvUpdateReq) error {
	go func() {
		if req.ID > 0 {
			// 更新链接，如果失败查看失败次数超过限值，则自动删除
			if req.Fail {
				take := l.svcCtx.TvModel.Take(req.ID)
				autoDelCount := l.svcCtx.SettingModel.Value("autoDelCount")
				to, _ := convertor.ToInt(autoDelCount)
				if to >= take.FailCount+1 {
					l.svcCtx.TvModel.Delete([]string{strconv.FormatUint(take.ID, 10)})
					return
				}
			}
			_ = l.svcCtx.TvModel.Update(model.Tv{ID: req.ID, Width: req.Width, Height: req.Height, Speed: req.Speed, DisplayName: req.DisplayName, Group: req.Group, Weight: req.Weight}, req.Fail)
		} else {
			if !req.Fail {
				_ = l.svcCtx.TvModel.Save(model.Tv{Name: req.Name, Url: req.Url, Width: req.Width, Height: req.Height, Speed: req.Speed, DisplayName: req.DisplayName, Group: req.Group}, req.Fail)
			}
		}
	}()

	return nil
}
