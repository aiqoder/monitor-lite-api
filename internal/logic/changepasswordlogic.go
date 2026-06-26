package logic

import (
	"context"
	"errors"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) error {
	if req.OldPassword == "" || req.NewPassword == "" {
		return errors.New("请填写原密码和新密码")
	}
	if len(req.NewPassword) < 4 {
		return errors.New("新密码至少 4 位")
	}
	if req.OldPassword == req.NewPassword {
		return errors.New("新密码不能与原密码相同")
	}

	current := l.svcCtx.SettingModel.Value("password")
	if req.OldPassword != current {
		return errors.New("原密码错误")
	}

	return l.svcCtx.SettingModel.UpdateValueByKey("password", req.NewPassword)
}
