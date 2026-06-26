package logic

import (
	"context"
	"strings"

	"github.com/aiqoder/monitor-lite-api/internal/pkg/prompt"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type RuleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRuleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RuleUpdateLogic {
	return &RuleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RuleUpdateLogic) RuleUpdate(req *types.RuleUpdateReq) error {
	groups := make([]prompt.RuleGroup, 0, len(req.Groups))
	for _, g := range req.Groups {
		name := strings.TrimSpace(g.Name)
		if name == "" {
			continue
		}
		channels := make([]string, 0, len(g.Channels))
		seen := make(map[string]struct{})
		for _, ch := range g.Channels {
			ch = strings.TrimSpace(ch)
			if ch == "" {
				continue
			}
			if _, ok := seen[ch]; ok {
				continue
			}
			seen[ch] = struct{}{}
			channels = append(channels, ch)
		}
		if len(channels) == 0 {
			continue
		}
		groups = append(groups, prompt.RuleGroup{
			Name:     name,
			Channels: channels,
		})
	}
	cfg := prompt.Config{Groups: groups}
	content, err := cfg.Marshal()
	if err != nil {
		return err
	}
	if err := l.svcCtx.SettingModel.SaveRule(content); err != nil {
		return err
	}
	l.svcCtx.PromptConfig = l.svcCtx.UpdatePrompt()
	return nil
}
