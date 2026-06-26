package logic

import (
	"context"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/aigroup"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type AiModelsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiModelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiModelsLogic {
	return &AiModelsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiModelsLogic) AiModels(req *types.AiModelsReq) (*types.AiModelsResp, error) {
	baseURL := req.BaseURL
	apiKey := req.ApiKey
	if baseURL == "" {
		baseURL = l.svcCtx.SettingModel.Value("aiBaseUrl")
	}
	if apiKey == "" {
		apiKey = l.svcCtx.SettingModel.Value("aiApiKey")
	}

	models, err := aigroup.ListModels(aigroup.AISettings{
		BaseURL: baseURL,
		APIKey:  apiKey,
	})
	if err != nil {
		return nil, err
	}
	return &types.AiModelsResp{Models: models}, nil
}
