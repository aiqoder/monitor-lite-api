package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SuperTvLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSuperTvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuperTvLogic {
	return &SuperTvLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuperTvLogic) SuperTv(req *types.SuperTvReq) (resp string, err error) {
	//读取规则文件
	return genTxt(l.svcCtx, true), nil
}

func (l *SuperTvLogic) SuperOutTv(req *types.SuperTvReq) (resp string, err error) {
	//读取规则文件
	return genTxt(l.svcCtx, false), nil
}
