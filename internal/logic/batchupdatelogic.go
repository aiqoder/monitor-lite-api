package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/model"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type BatchUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateLogic {
	return &BatchUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchUpdateLogic) BatchUpdate(req *types.BatchTvUpdateReq) error {
	for i := range req.Tvs {
		tv := req.Tvs[i]
		_ = l.svcCtx.TvModel.Save(model.Tv{Url: tv.Url, Name: tv.Name}, tv.Fail)
	}

	return nil
}
