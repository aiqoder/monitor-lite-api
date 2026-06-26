package logic

import (
	"context"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type BatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLogic {
	return &BatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteLogic) BatchDelete(req *types.BatchTvDelReq) error {
	l.svcCtx.TvModel.Delete(req.Ids)
	return nil
}
