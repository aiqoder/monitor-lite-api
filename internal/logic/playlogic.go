package logic

import (
	"context"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type PlayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlayLogic {
	return &PlayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlayLogic) Play(req *types.PlayReq) error {
	// todo: add your logic here and delete this line

	return nil
}
