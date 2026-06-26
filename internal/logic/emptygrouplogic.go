package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

type EmptyGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmptyGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmptyGroupLogic {
	return &EmptyGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmptyGroupLogic) EmptyGroup() error {
	return l.svcCtx.TvModel.EmptyGroup()
}
