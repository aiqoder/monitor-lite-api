package router

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	epg "github.com/aiqoder/monitor-lite-api/internal/handler/epg"
	"github.com/aiqoder/monitor-lite-api/internal/handler"
	proxy "github.com/aiqoder/monitor-lite-api/internal/handler/proxy"
	selfout "github.com/aiqoder/monitor-lite-api/internal/handler/selfout"
	subscriber "github.com/aiqoder/monitor-lite-api/internal/handler/subscriber"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/pkg/common/jwt"
	"github.com/aiqoder/monitor-lite-api/pkg/common/mime"
	"github.com/aiqoder/monitor-lite-api/pkg/common/tools"
	"github.com/aiqoder/monitor-lite-api/logo"
	"github.com/aiqoder/monitor-lite-api/ui"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Register(r *gin.Engine, ctx *svc.ServiceContext) {
	registerAPI(r, ctx)
	registerStatic(r)
	registerLogo(r)
	registerWebSocket(r, ctx)
}

func registerAPI(r *gin.Engine, ctx *svc.ServiceContext) {
	tv := r.Group("/v1/tv")
	{
		tv.GET("/identify", handler.IdentifyHandler(ctx))
		tv.GET("/json", handler.SearchHandler(ctx))
		tv.GET("/page", handler.PageHandler(ctx))
		tv.GET("/pix", handler.PixHandler(ctx))
		tv.GET("/super", handler.SuperTvHandler(ctx))
		tv.GET("/w/:path", handler.SuperWatchTvHandler(ctx))
		tv.POST("/update", handler.UpdateHandler(ctx))
		tv.POST("/check", handler.CheckHandler(ctx))
		tv.POST("/batchupdate", handler.BatchUpdateHandler(ctx))
		tv.POST("/batchdelete", handler.BatchDeleteHandler(ctx))
		tv.POST("/deleteLoseEfficacy", handler.DeleteLoseEfficacyHandler(ctx))
		tv.POST("/emptyGroup", handler.EmptyGroupHandler(ctx))
		tv.POST("/updateGroup", handler.UpdateGroupHandler(ctx))
		tv.POST("/checkAll", handler.CheckAllHandler(ctx))
		tv.GET("/rule/get", handler.RuleGetHandler(ctx))
		tv.POST("/rule/update", handler.RuleUpdateHandler(ctx))
		tv.GET("/sync/mysql", handler.SyncTVDataHandler(ctx))
		tv.GET("/play", handler.PlayHandler(ctx))
		tv.GET("/tips", ctx.WSHub.HandleWebSocketGin)
	}

	setting := r.Group("/v1/setting")
	{
		setting.GET("/find", handler.FindSettingHandler(ctx))
		setting.POST("/update", handler.UpdateSettingHandler(ctx))
		setting.POST("/changePassword", handler.ChangePasswordHandler(ctx))
		setting.GET("/aiModels", handler.AiModelsHandler(ctx))
	}

	proxyGroup := r.Group("/v1/proxy")
	{
		proxyGroup.GET("/fengshows", proxy.FengShowsHandler(ctx))
		proxyGroup.GET("/ptbtv", proxy.PtbtvHandler(ctx))
		proxyGroup.GET("/iqilu", proxy.IQiLuHandler(ctx))
	}

	sub := r.Group("/v1/subscriber")
	{
		sub.GET("/subscribers", subscriber.SubscriberListHandler(ctx))
		sub.POST("/update", subscriber.SubscriberUpdateHandler(ctx))
		sub.GET("/delete/:id", subscriber.SubscriberDeleteHandler(ctx))
		sub.GET("/grab/:id", subscriber.SubscriberGrabHandler(ctx))
	}

	epgGroup := r.Group("/v1/epg")
	{
		epgGroup.GET("/epgList", epg.EpgListHandler(ctx))
		epgGroup.GET("/diyp", epg.EpgSearchHandler(ctx))
		epgGroup.GET("/collect", epg.EpgCollectHandler(ctx))
	}

	r.GET("/cus/:path", handler.CustomOutHandler(ctx))

	selfoutGroup := r.Group("/v1/selfout")
	{
		selfoutGroup.POST("/addSelfout", selfout.AddSelfoutHandler(ctx))
		selfoutGroup.POST("/updateSelfout", selfout.UpdateSelfoutHandler(ctx))
		selfoutGroup.DELETE("/delSelfout", selfout.DelSelfoutHandler(ctx))
		selfoutGroup.GET("/getSelfoutById", selfout.GetSelfoutHandler(ctx))
		selfoutGroup.GET("/searchSelfout", selfout.SearchSelfoutHandler(ctx))
	}
}

func registerStatic(r *gin.Engine) {
	serveAdmin := func(c *gin.Context) {
		path := strings.Replace(c.Request.RequestURI, "/admin", "", 1)
		getPath := func(p string) string {
			prefix := "dist%s"
			if p == "/" || filepath.Ext(p) == "" {
				return fmt.Sprintf(prefix, "/index.html")
			}
			return fmt.Sprintf(prefix, p)
		}

		filePath := getPath(path)
		bytes, err := ui.DIST.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.Header("content-type", mime.TypeByExtension(filepath.Ext(filePath)))
		c.Header("Cache-Control", "public,max-age=86400")
		c.Writer.Write(bytes)
	}

	r.GET("/admin", serveAdmin)
	r.GET("/admin/*filepath", serveAdmin)
}

func registerLogo(r *gin.Engine) {
	r.GET("/logo/*path", func(c *gin.Context) {
		writer := func(file []byte) {
			c.Header("content-type", "image/*")
			c.Writer.Write(file)
		}

		path, _ := url.PathUnescape(c.Request.RequestURI)
		wordMap := map[string]string{
			"⁺": "+", "丨": "", "《": "", "》": "", "·": "",
			"，": "", "：": "", "哒啵": "哒波",
		}
		for key, val := range wordMap {
			path = strings.Replace(path, key, val, -1)
		}

		joinPath := strings.Replace(path, "logo", "images", -1)
		_, fileName := filepath.Split(joinPath)
		imageName := strings.Replace(fileName, ".png", "", -1)
		file, err := logo.Logo.ReadFile(joinPath[1:])

		file, err = logo.Logo.ReadFile(fmt.Sprintf("images/%s.png", imageName))
		if err == nil {
			writer(file)
			return
		}

		if !tools.IsContainsEnglish(imageName) {
			suffix := []string{"频道", "台", "套", "剧"}
			if err != nil {
				for i := 0; i < len(suffix); i++ {
					file, _ = logo.Logo.ReadFile(fmt.Sprintf("images/%s%s.png", imageName, suffix[i]))
					if file != nil {
						writer(file)
						return
					}
				}
				for i := 0; i < len(suffix); i++ {
					before, found := strings.CutSuffix(imageName, suffix[i])
					if found {
						file, _ = logo.Logo.ReadFile(fmt.Sprintf("images/%s.png", before))
						if file != nil {
							writer(file)
							return
						}
					}
				}
			}
			writer(file)
			return
		}

		file, err = logo.Logo.ReadFile(fmt.Sprintf("images/%s.png", strings.ToLower(imageName)))
		if err == nil {
			writer(file)
			return
		}

		file, err = logo.Logo.ReadFile(fmt.Sprintf("images/%s.png", strings.ToUpper(imageName)))
		if err == nil {
			writer(file)
			return
		}

		file, err = logo.Logo.ReadFile(fmt.Sprintf("images/%s.png", strutil.CamelCase(imageName)))
		if err == nil {
			writer(file)
			return
		}

		fmt.Println(c.Request.RequestURI)
		file, err = logo.Logo.ReadFile("images/橙子.png")
		if err != nil {
			writer(file)
		}
	})
}

func registerWebSocket(r *gin.Engine, ctx *svc.ServiceContext) {
	r.GET("/v1/video/play", func(c *gin.Context) {
		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Errorf("Error upgrading to WebSocket: %v", err)
			return
		}
		query := c.Request.URL.Query()
		playUrl := query["url"][0]

		reader, writer := io.Pipe()

		cmd := ffmpeg.Input(playUrl).
			Output("pipe:", ffmpeg.KwArgs{
				"c:v": "copy",
				"c:a": "aac",
				"f":   "flv",
			}).WithOutput(writer)

		_, cancelPlayer := context.WithCancel(cmd.Context)

		stop := func(err error) {
			log.Errorf("[ffmpeg 播放失败]：%v，[URL：%s]", err, playUrl)
			cancelPlayer()
			_ = conn.Close()
			_ = writer.Close()
			_ = reader.Close()
		}

		go func() {
			err = cmd.Run()
			if err != nil {
				stop(err)
			}
		}()

		lastSendTime := time.Now().UnixMilli()
		go func() {
			for {
				mt, msg, rErr := conn.ReadMessage()
				if rErr != nil {
					stop(rErr)
					break
				}
				if mt == websocket.TextMessage {
					lastSendTime, _ = convertor.ToInt(string(msg))
				}
			}
		}()

		go func() {
			for {
				timer := time.NewTimer(time.Second)
				<-timer.C
				if time.Now().UnixMilli()-lastSendTime > 1000*10 {
					stop(errors.New("视频播放已结束"))
					break
				}
			}
		}()

		bufReader := bufio.NewReader(reader)
		for {
			bytes, err := bufReader.ReadBytes('\n')
			if err == io.EOF {
				stop(err)
				break
			}
			err = conn.WriteMessage(websocket.BinaryMessage, bytes)
			if err != nil {
				stop(err)
				break
			}
		}
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		if tools.StrInContains(uri, []string{
			"/cus", "/video/play", "/v1/tv/w/", "/v1/tv/tips",
			"/v1/tv/identify", "/admin", "/epg/diyp", "/epg/collect",
			"/v1/proxy", "/logo",
		}) {
			c.Next()
			return
		}

		auth := c.GetHeader("session-id")
		if _, err := jwt.ParseToken(auth); err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
