package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/pkg/common/ffmpeg"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
	"github.com/aiqoder/monitor-lite-api/model"
)

type CheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLogic {
	return &CheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckLogic) Check(req *types.TvCheckReq) error {
	info, err := ffmpeg.GetOnlineVideoInfo(req.Url)
	if err != nil {
		return err
	}
	_ = l.svcCtx.TvModel.Save(model.Tv{Name: req.Name, Url: req.Url, Speed: info.Speed, Width: info.Width, Height: info.Height}, false)
	return nil
}
