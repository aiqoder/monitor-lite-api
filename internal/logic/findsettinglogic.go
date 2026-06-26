package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type FindSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindSettingLogic {
	return &FindSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindSettingLogic) FindSetting(req *types.FindReq) (resp *types.FindResp, err error) {
	list, err := l.svcCtx.SettingModel.TableList()
	var settings []types.Setting
	for i := 0; i < len(list); i++ {
		settings = append(settings, types.Setting{
			ID:    list[i].ID,
			Key:   list[i].Key,
			Value: list[i].Value,
			Type:  list[i].Type,
		})
	}

	return &types.FindResp{
		Settings: settings,
	}, nil
}
