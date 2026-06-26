package selfout

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/model"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type UpdateSelfoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSelfoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSelfoutLogic {
	return &UpdateSelfoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSelfoutLogic) UpdateSelfout(req *types.UpdateSelfoutReq) (resp *types.UpdateSelfoutResp, err error) {
	err = l.svcCtx.SelfoutModel.Save(&model.Selfout{
		Id:     req.Id,
		Width:  req.Width,
		Height: req.Height,
		Speed:  req.Speed,
		Key:    req.Key,
	})

	return
}
