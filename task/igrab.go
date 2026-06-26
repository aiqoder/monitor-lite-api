package task

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"strings"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/tools"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/model"
)

func Grab(ctx *svc.ServiceContext, item model.Subscriber) {
	source := tools.ParserIpTvSource(item.Url)
	var tvs []model.Tv
	for i := 0; i < len(source); i++ {
		s := source[i]
		t := time.Now()
		url := strings.Replace(s.Url, "\r", "", -1)
		tvs = append(tvs, model.Tv{
			Name:       s.Name,
			Url:        url,
			UpdateTime: &t,
		})

	}

	err2 := ctx.TvModel.BatchSave(tvs)

	if err2 != nil {
		log.Error(err2)
	}

	t := time.Now()
	err := ctx.SubscriberModel.Update(model.Subscriber{
		ID:        item.ID,
		Name:      item.Name,
		Url:       item.Url,
		Count:     uint64(len(source)),
		CheckTime: &t,
	})

	if err != nil {
		log.Error(err)
	}
}

// 自动抓取订阅
func Igrab(ctx *svc.ServiceContext) func() {
	return func() {
		for {
			timer := time.NewTimer(time.Second * 30)
			<-timer.C
			now, _ := time.Parse("15:04", time.Now().Format("15:04"))
			subscriberTime := ctx.SettingModel.Value("subscriberTime")
			parse, err := time.Parse("15:04", subscriberTime)
			if err != nil {
				log.Error(err)
				continue
			}

			diff := now.Unix() - parse.Unix()

			if diff >= 0 && diff <= 60 {
				list, err := ctx.SubscriberModel.List()
				if err != nil {
					log.Error(err)
					return
				}
				for _, item := range list {
					Grab(ctx, item)
				}
			}
		}
	}
}
