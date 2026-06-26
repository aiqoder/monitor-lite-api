package subscriber

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/model"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SubscriberUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscriberUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriberUpdateLogic {
	return &SubscriberUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriberUpdateLogic) SubscriberUpdate(req *types.Subscriber) (resp string, err error) {
	err = l.svcCtx.SubscriberModel.Update(model.Subscriber{
		ID:   req.ID,
		Name: req.Name,
		Url:  req.Url,
	})

	if err != nil {
		return "保存失败", err
	}

	return "保存成功", nil
}
