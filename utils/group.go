package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/aigroup"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/model"
)

func AISettingsFromCtx(ctx *svc.ServiceContext) aigroup.AISettings {
	return aigroup.AISettings{
		BaseURL: ctx.SettingModel.Value("aiBaseUrl"),
		APIKey:  ctx.SettingModel.Value("aiApiKey"),
		Model:   ctx.SettingModel.Value("aiModel"),
	}
}

func UpdateGroupByAI(ctx *svc.ServiceContext, tvList []model.Tv) error {
	settings := AISettingsFromCtx(ctx)
	client := aigroup.NewClient(settings)
	if !client.Enabled() {
		return fmt.Errorf("请先在系统设置中配置 AI API Key")
	}

	rawNames := make([]string, 0, len(tvList))
	indexMap := make([]int, 0, len(tvList))
	for i := range tvList {
		tv := tvList[i]
		if validator.IsNumberStr(tv.Name) {
			continue
		}
		if len(strings.TrimSpace(tv.Name)) == 0 {
			continue
		}
		indexMap = append(indexMap, i)
		rawNames = append(rawNames, tv.Name)
	}

	if len(rawNames) == 0 {
		return nil
	}

	results, err := aigroup.ClassifyAll(settings, ctx.PromptConfig, rawNames)
	if err != nil {
		return err
	}

	for j, tvIdx := range indexMap {
		r, ok := results[j]
		if !ok || r.Group == "" {
			continue
		}
		tv := tvList[tvIdx]
		displayName := r.DisplayName
		if displayName == "" {
			displayName = tv.Name
		}
		err := ctx.TvModel.UpdatesByName(model.Tv{
			Name:        tv.Name,
			DisplayName: displayName,
			Group:       r.Group,
		})
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func UpdateGroup(ctx *svc.ServiceContext, tvList []*model.Tv) {
	start := time.Now()
	_ = UpdateGroupByAI(ctx, modelTvsFromPtrs(tvList))
	fmt.Println("AI 分组完成 time", time.Since(start).Milliseconds())
}

func modelTvsFromPtrs(list []*model.Tv) []model.Tv {
	out := make([]model.Tv, 0, len(list))
	for _, tv := range list {
		if tv != nil {
			out = append(out, *tv)
		}
	}
	return out
}
