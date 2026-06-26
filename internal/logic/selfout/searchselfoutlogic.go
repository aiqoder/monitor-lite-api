package selfout

import (
	"context"
	"strconv"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SearchSelfoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchSelfoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSelfoutLogic {
	return &SearchSelfoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchSelfoutLogic) SearchSelfout(req *types.SearchSelfoutReq) (resp *types.SearchSelfoutResp, err error) {
	selfout, count, err := l.svcCtx.SelfoutModel.ListPage(
		map[string]string{
			"current": strconv.FormatInt(req.Current, 10),
			"size":    strconv.FormatInt(req.Size, 10),
			"id":      strconv.FormatInt(req.Id, 10),
			"width":   strconv.FormatInt(req.Width, 10),
			"height":  strconv.FormatInt(req.Height, 10),
			"speed":   strconv.FormatInt(req.Speed, 10),
			"key":     req.Key,
		})

	if err != nil {
		return nil, err
	}

	return &types.SearchSelfoutResp{
		Total:   count,
		Size:    req.Size,
		Current: req.Current,
		Data:    selfout,
	}, nil
}
