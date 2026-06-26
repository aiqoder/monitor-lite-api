package selfout

import (
	"context"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type GetSelfoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelfoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfoutLogic {
	return &GetSelfoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelfoutLogic) GetSelfout(req *types.GetSelfoutByIdReq) (resp *types.GetSelfoutByIdResp, err error) {
	one, err := l.svcCtx.SelfoutModel.FirstOne(req.Id)

	if err != nil {
		return nil, err
	}
	return &types.GetSelfoutByIdResp{
		Width:  one.Width,
		Height: one.Height,
		Speed:  one.Speed,
		Key:    one.Key,
	}, nil
}
