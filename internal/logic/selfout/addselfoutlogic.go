package selfout

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/model"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type AddSelfoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSelfoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSelfoutLogic {
	return &AddSelfoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSelfoutLogic) AddSelfout(req *types.AddSelfoutReq) (resp *types.AddSelfoutResp, err error) {
	err = l.svcCtx.SelfoutModel.Save(&model.Selfout{
		Width:  req.Width,
		Height: req.Height,
		Speed:  req.Speed,
		Key:    req.Key,
	})

	if err != nil {
		return nil, err
	}

	return &types.AddSelfoutResp{}, nil
}
