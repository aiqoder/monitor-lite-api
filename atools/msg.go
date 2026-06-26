package atools

import (
	"fmt"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

// 用于全局msg消息通知
func WsMsg(svcCtx *svc.ServiceContext, title, content string) {
	msg := fmt.Sprintf("%s#%s", title, content)
	svcCtx.WSHub.Broadcast <- []byte(msg)
}
