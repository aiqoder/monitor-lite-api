package subscriber

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/task"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SubscriberGrabLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscriberGrabLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriberGrabLogic {
	return &SubscriberGrabLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriberGrabLogic) SubscriberGrab(req *types.SubscriberDeleteReq) (resp string, err error) {
	item, err := l.svcCtx.SubscriberModel.GetById(req.ID)

	if err != nil {
		return "", err
	}
	task.Grab(l.svcCtx, item)
	return "抓取完成", nil
}
