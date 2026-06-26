package subscriber

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"context"
	"time"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SubscriberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscriberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriberListLogic {
	return &SubscriberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscriberListLogic) SubscriberList() (resp []types.Subscriber, err error) {
	list, err := l.svcCtx.SubscriberModel.List()

	var result []types.Subscriber

	if err != nil {
		log.Error(err)
		return result, err
	}

	for i := 0; i < len(list); i++ {
		l := list[i]
		s := types.Subscriber{
			ID:    l.ID,
			Name:  l.Name,
			Url:   l.Url,
			Count: l.Count,
		}

		if l.CheckTime != nil {
			s.CheckTime = l.CheckTime.Format(time.DateTime)
		}
		result = append(result, s)
	}
	return result, err
}
