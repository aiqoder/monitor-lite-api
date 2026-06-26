package logic

import (
	"context"
	"errors"
	"github.com/duke-git/lancet/v2/cryptor"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/jwt"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type IdentifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentifyLogic {
	return &IdentifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentifyLogic) Identify(req *types.IdentifyReq) (resp *types.IdentifyResp, err error) {
	username := l.svcCtx.SettingModel.Value("username")
	password := l.svcCtx.SettingModel.Value("password")
	if req.Username != username || req.Password != password {
		return nil, errors.New("授权失败")
	}

	token, _ := jwt.GenerateToken(jwt.Claims{
		Username: username,
		UserId:   0,
		AppId:    password,
		MapClaims: jwt2.MapClaims{
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
			"Issuer": "YiGeChengZi",
		},
	})

	return &types.IdentifyResp{Password: cryptor.Md5String(password), Token: token}, nil
}
