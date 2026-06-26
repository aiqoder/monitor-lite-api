package svc

import (
	"io"
	"os"
	"time"

	"gorm.io/gorm/logger"
	"github.com/aiqoder/monitor-lite-api/pkg/common/db"
	"github.com/aiqoder/monitor-lite-api/internal/config"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/cache"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/prompt"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/rulecache"
	"github.com/aiqoder/monitor-lite-api/model"
	"github.com/aiqoder/monitor-lite-api/websocket"

	"github.com/duke-git/lancet/v2/fileutil"
)

type ServiceContext struct {
	Config          config.Config
	TvModel         model.TvModel
	SettingModel    model.SettingModel
	EpgModel        model.EpgModel
	SubscriberModel model.SubscriberModel
	SelfoutModel    model.SelfoutModel
	Cache           *cache.Cache
	WSHub           *websocket.Hub
	PromptConfig    prompt.Config
	UpdatePrompt    func() prompt.Config
}

func seedRuleContent(settingModel model.SettingModel) {
	if settingModel.RuleContent() != "" {
		return
	}
	// 仅写入空配置；AI 分类在 channels 为空时会回退到 defaultgroup 内置数据
	settingModel.EnsureRule(`{"groups":[]}`)
}

func NewServiceContext(c config.Config) *ServiceContext {
	appCache, err := cache.New(time.Second*2, 10000)
	if err != nil {
		panic(err)
	}

	if fileutil.IsExist("sqlite.db") && !fileutil.IsExist("etc/sqlite.db") {
		distFile, err := os.Create("etc/sqlite.db")
		if err != nil {
			panic(err)
		}

		sourceFile, err := os.Open("sqlite.db")
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(distFile, sourceFile)
		if err != nil {
			panic(err)
		}
	}

	dbCfg := db.NewDB(db.NewSqlite("etc/sqlite.db"))

	if c.Mode == "dev" {
		dbCfg.DB.Logger = logger.Default.LogMode(logger.Info)
	} else {
		dbCfg.DB.Logger = logger.Default.LogMode(logger.Error)
	}

	_ = dbCfg.DB.AutoMigrate(&model.Tv{}, &model.Epg{}, model.Selfout{}, &model.Setting{}, &model.Subscriber{})

	wsHub := websocket.NewHub()
	go wsHub.Run()

	settingModel := model.NewSettingModel(dbCfg)
	seedRuleContent(settingModel)

	promptCfg := prompt.Config{}
	updatePrompt := func() prompt.Config {
		content := settingModel.RuleContent()
		cfg, err := prompt.Parse(content)
		if err != nil {
			log.Error(err)
			cfg = prompt.DefaultConfig()
		}
		promptCfg = cfg

		rulecache.SetGroups(cfg.GroupMap())

		return promptCfg
	}

	updatePrompt()

	return &ServiceContext{
		Config:          c,
		Cache:           appCache,
		SettingModel:    settingModel,
		TvModel:         model.NewTvModel(dbCfg),
		SubscriberModel: model.NewSubscriberModel(dbCfg),
		EpgModel:        model.NewEpgModel(dbCfg),
		SelfoutModel:    model.NewSelfoutModel(dbCfg),
		WSHub:           wsHub,
		PromptConfig:    promptCfg,
		UpdatePrompt:    updatePrompt,
	}
}
