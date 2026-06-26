package subscriber

import (
	"context"
	"errors"
	"github.com/aiqoder/monitor-lite-api/model"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SubscriberDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscriberDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriberDeleteLogic {
	return &SubscriberDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriberDeleteLogic) SubscriberDelete(req *types.SubscriberDeleteReq) (resp string, err error) {
	if req.ID == 0 {
		return "", errors.New("id is required")
	}

	l.svcCtx.SubscriberModel.Delete(model.Subscriber{
		ID: req.ID,
	})

	return "删除成功", nil
}
