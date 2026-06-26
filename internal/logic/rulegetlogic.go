package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/prompt"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type RuleGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRuleGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RuleGetLogic {
	return &RuleGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RuleGetLogic) RuleGet() (resp types.RuleGetResp, err error) {
	cfg, err := prompt.Parse(l.svcCtx.SettingModel.RuleContent())
	if err != nil {
		return types.RuleGetResp{}, err
	}
	groups := make([]types.RuleGroup, 0, len(cfg.Groups))
	for _, g := range cfg.Groups {
		groups = append(groups, types.RuleGroup{
			Name:     g.Name,
			Channels: g.Channels,
		})
	}
	return types.RuleGetResp{Groups: groups}, nil
}
