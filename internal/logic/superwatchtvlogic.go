package logic

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/prompt"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SuperWatchTvLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func getGroup(str string) []string {
	cfg, err := prompt.Parse(str)
	if err == nil && len(cfg.Groups) > 0 {
		return cfg.GroupNames()
	}
	re, _ := regexp.Compile(" (.*?):")
	groupStr := strings.Split(str, "name:")
	gr := re.FindAllString(groupStr[0], 100)
	var groupNames []string
	for i := 0; i < len(gr); i++ {
		g := gr[i]
		groupNames = append(groupNames, strings.TrimSpace(strings.TrimSuffix(g, ":")))
	}
	return groupNames
}

func NewSuperWatchTvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuperWatchTvLogic {
	return &SuperWatchTvLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuperWatchTvLogic) SuperWatchTv(req *types.SuperWatchTvReq) (resp string, err error) {
	path := req.Path
	plusKey := l.svcCtx.SettingModel.Value("plusKey")

	if strings.HasSuffix(req.Path, ".m3u") {
		path, _ = strings.CutSuffix(req.Path, ".m3u")
	}

	if path != plusKey {
		return "", errors.New("非法访问，这里毛都没有")
	}

	if strings.HasSuffix(req.Path, ".m3u") {
		return genM3u(l.svcCtx), nil
	}

	if path == "" {
		return "", errors.New("非法访问，这里毛都没有")
	}

	return genTxt(l.svcCtx, true), nil
}

func (l *SuperWatchTvLogic) SuperOutWatchTv(req *types.SuperWatchTvReq) (resp string, err error) {
	path := req.Path
	plusKey := l.svcCtx.SettingModel.Value("plusKey")

	if strings.HasSuffix(req.Path, ".m3u") {
		path, _ = strings.CutSuffix(req.Path, ".m3u")
	}

	if path != plusKey {
		return "", errors.New("非法访问，这里毛都没有")
	}

	if strings.HasSuffix(req.Path, ".m3u") {
		return genM3u(l.svcCtx), nil
	}

	if path == "" {
		return "", errors.New("非法访问，这里毛都没有")
	}

	return genTxt(l.svcCtx, false), nil
}
