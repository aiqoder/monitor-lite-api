package proxy

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type FengShowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFengShowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FengShowsLogic {
	return &FengShowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FengShowsLogic) FengShows(req *types.FengShowReq) (resp string, err error) {
	var request = resty.New()
	//7c96b084-60e1-40a9-89c5-682b994fb680  资讯台
	//f7f48462-9b13-485b-8101-7b54716411ec  中文台
	//15e02d92-1698-416c-af2f-3e9a872b4d78  深圳旁边台
	vid := map[string]string{
		"zx": "7c96b084-60e1-40a9-89c5-682b994fb680",
		"zw": "f7f48462-9b13-485b-8101-7b54716411ec",
		"xg": "15e02d92-1698-416c-af2f-3e9a872b4d78",
	}
	bstrURL := fmt.Sprintf("https://m.fengshows.com/api/v3/hub/live/auth-url?live_id=%s&live_qa=hd", vid[req.T])
	resultMap := struct {
		Data struct {
			LiveUrl string `json:"live_url"`
		} `json:"data"`
		Status  string
		Message string
	}{}
	_, err = request.R().SetResult(&resultMap).Get(bstrURL)

	if err != nil {
		return "", err
	}

	live := resultMap.Data.LiveUrl

	return live, nil
}
