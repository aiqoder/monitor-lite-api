package logic

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

type PixLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPixLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PixLogic {
	return &PixLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PixLogic) Pix() (resp []string, err error) {
	pixs := l.svcCtx.TvModel.Pixs()

	if pixs == nil {
		return nil, fmt.Errorf("not found pixs")
	}

	var result []string

	for _, pix := range pixs {
		result = append(result, fmt.Sprintf("%dX%d", pix.Width, pix.Height))
	}
	unique := slice.Unique(result)
	return unique, nil
}
