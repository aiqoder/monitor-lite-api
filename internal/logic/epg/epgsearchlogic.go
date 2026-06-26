package epg

import (
	"context"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type EpgSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEpgSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EpgSearchLogic {
	return &EpgSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EpgSearchLogic) EpgSearch(req *types.EpgSearchReq) (resp types.EpgSearchReply, err error) {
	// DIYP 判断条件
	if len(req.Date) > 0 && len(req.Ch) > 0 {
		epgs := l.svcCtx.EpgModel.ChannelEpgByDate(req.Ch, req.Date)

		var resultEpgData []types.EpgData

		for i := 0; i < len(epgs); i++ {
			epg := epgs[i]
			resultEpgData = append(resultEpgData, types.EpgData{
				Channel: epg.Channel,
				Start:   epg.Start[:5],
				Stop:    epg.Stop[:5],
				Title:   epg.Title,
				Desc:    epg.Desc,
				Date:    epg.Date,
			})
		}

		return types.EpgSearchReply{
			ChannelName: req.Ch,
			Date:        req.Date,
			EpgData:     resultEpgData,
		}, nil
	}

	return types.EpgSearchReply{}, nil
}
