package selfout

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/model"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type DelSelfoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelSelfoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelSelfoutLogic {
	return &DelSelfoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelSelfoutLogic) DelSelfout(req *types.DelSelfoutReq) (resp *types.DelSelfoutResp, err error) {
	err = l.svcCtx.SelfoutModel.Del(&model.Selfout{Id: req.Id})

	return
}
