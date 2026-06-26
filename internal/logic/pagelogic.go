package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageLogic) Page(req *types.TvPageReq) (resp types.TvListPageResp, err error) {
	var result []types.TvListResp

	page, total, err := l.svcCtx.TvModel.TableListPage(map[string]string{
		"current":     strconv.FormatInt(req.Current, 10),
		"size":        strconv.FormatInt(req.Size, 10),
		"name":        req.Name,
		"displayName": req.DisplayName,
		"group":       req.Group,
		"url":         req.Url,
		"width":       req.Width,
		"height":      req.Height,
	})

	for i := 0; i < len(page); i++ {
		tv := page[i]
		result = append(result, types.TvListResp{
			ID:          tv.ID,
			Name:        tv.Name,
			Url:         tv.Url,
			Width:       tv.Width,
			Height:      tv.Height,
			FailCount:   uint64(tv.FailCount),
			UpdateTime:  tv.UpdateTime.Format(time.DateTime),
			Speed:       tv.Speed,
			Group:       tv.Group,
			DisplayName: tv.DisplayName,
			Weight:      tv.Weight,
		})
	}

	return types.TvListPageResp{
		Records: result,
		Total:   total,
		Current: req.Current,
		Size:    req.Size,
	}, err
}
