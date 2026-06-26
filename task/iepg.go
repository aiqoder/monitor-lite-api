package task

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
)

func IEpg(ctx *svc.ServiceContext) func() {
	return func() {
		if ctx.SettingModel.Value("autoEpg") == "0" {
			log.Error("未开启EPG")
			return
		}
		ctx.EpgModel.UpdateXml2db()
	}
}
