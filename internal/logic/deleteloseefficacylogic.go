package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

type DeleteLoseEfficacyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLoseEfficacyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoseEfficacyLogic {
	return &DeleteLoseEfficacyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLoseEfficacyLogic) DeleteLoseEfficacy() error {
	return l.svcCtx.TvModel.DeleteAllFail()
}
