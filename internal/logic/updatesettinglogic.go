package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/model"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type UpdateSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSettingLogic {
	return &UpdateSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSettingLogic) UpdateSetting(req *types.UpdateSettingReq) (resp *types.UpdateSettingResp, err error) {
	ss := req.Settings

	for i := 0; i < len(ss); i++ {
		s := ss[i]
		_ = l.svcCtx.SettingModel.Update(model.Setting{
			ID:    s.ID,
			Key:   s.Key,
			Value: s.Value,
			Type:  s.Type,
		})
	}

	return &types.UpdateSettingResp{}, nil
}
