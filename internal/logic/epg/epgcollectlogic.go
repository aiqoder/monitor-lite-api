package epg

import (
	"context"
	"errors"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

type EpgCollectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEpgCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EpgCollectLogic {
	return &EpgCollectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var epgCollectFlag = false

func (l *EpgCollectLogic) EpgCollect() error {
	if epgCollectFlag {
		return errors.New("后台正在采集中，请勿重复下发命令")
	}
	go func() {
		epgCollectFlag = true
		l.svcCtx.EpgModel.UpdateXml2db()
		epgCollectFlag = false
	}()
	return nil
}
