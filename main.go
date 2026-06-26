package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/aiqoder/monitor-lite-api/pkg/common/fileutils"
	"github.com/aiqoder/monitor-lite-api/pkg/common/tools"
	"github.com/aiqoder/monitor-lite-api/etcexample"
	"github.com/aiqoder/monitor-lite-api/internal/config"
	"github.com/aiqoder/monitor-lite-api/internal/router"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/version"
	"github.com/aiqoder/monitor-lite-api/task"
)

var confPath = "etc/tv.yaml"

func main() {
	if !tools.IsCommandAvailable("ffmpeg") {
		panic("请先安装ffmpeg...")
	}

	cfile, _ := etcexample.EtcExample.ReadFile(confPath)
	fileutils.CreateFile(confPath, cfile)

	var configFile = flag.String("f", confPath, "the config file")
	flag.Parse()

	var c config.Config
	config.MustLoad(*configFile, &c)

	if c.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := svc.NewServiceContext(c)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(router.AuthMiddleware())
	router.Register(engine, ctx)

	fmt.Printf("Starting server at %s:%d (version %s)...\n", c.Host, c.Port, version.Version)
	task.Start(ctx)

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	if err := engine.Run(addr); err != nil {
		panic(err)
	}
}
