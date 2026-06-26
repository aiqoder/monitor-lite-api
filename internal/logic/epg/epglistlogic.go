package epg

import (
	"context"
	"errors"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type EpgListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEpgListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EpgListLogic {
	return &EpgListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EpgListLogic) EpgList(req *types.EpgListReq) (resp []types.EpgData, err error) {
	list := l.svcCtx.EpgModel.List(req.Channel)

	var result []types.EpgData

	if list == nil {
		return result, errors.New("没有查到节目信息")
	}

	for i := 0; i < len(list); i++ {
		l := list[i]
		result = append(result, types.EpgData{
			Channel: l.Channel,
			Start:   l.Start,
			Stop:    l.Stop,
			Title:   l.Title,
			Desc:    l.Desc,
			Date:    l.Date,
		})
	}
	return result, nil
}
