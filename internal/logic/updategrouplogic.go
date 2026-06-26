package logic

import (
	"context"
	"errors"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/task"
)

type UpdateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupLogic {
	return &UpdateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var updatedFlag = false

func (l *UpdateGroupLogic) UpdateGroup() error {
	if updatedFlag {
		return errors.New("后台正在更新分组，请勿重复下发指令！")
	}
	updatedFlag = true
	go task.UpdateGroup(l.svcCtx, "manual")()
	updatedFlag = false
	return nil
}
