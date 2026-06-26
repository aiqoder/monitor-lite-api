package logic

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"

	"github.com/aiqoder/monitor-lite-api/internal/svc"
	"github.com/aiqoder/monitor-lite-api/internal/types"
)

type SyncTVDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncTVDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncTVDataLogic {
	return &SyncTVDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncTVDataLogic) SyncTVData(req *types.SyncTVDataReq) error {
	host := "localhost"
	port := "3306"
	user := req.User
	password := req.Password
	database := req.Database
	table := req.Table
	insertJson := req.InsertJson
	insertBeforeClear := req.InsertBeforeClear

	if req.Host != "" {
		host = req.Host
	}

	if req.Port != "" {
		host = req.Host
	}

	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(source))

	if err != nil {
		log.Error(err)
		return err
	}

	if insertBeforeClear {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table))
	}

	list, err := l.svcCtx.TvModel.TableList()

	if err != nil {
		log.Error(err)
		return err
	}

	m := map[string]string{}

	err = json.Unmarshal([]byte(insertJson), &m)

	if err != nil {
		log.Error(err)
		return err
	}

	values := ""

	name := m["name"]
	url := m["url"]
	category := m["category"]

	for _, v := range list {
		if len(v.Group) == 0 {
			continue
		}
		if v.FailCount > 0 {
			continue
		}
		if v.Width == 0 {
			continue
		}
		split := strings.Split(v.Group, ",")
		for _, g := range split {
			values += fmt.Sprintf("('%s','%s','%s'),", v.Name, v.Url, g)
		}
	}
	before, found := strings.CutSuffix(values, ",")
	sql := fmt.Sprintf("INSERT INTO `%s` (%s,%s,%s) VALUES %s;", table, name, url, category, before)
	if found {
		err := db.Exec(sql).Error
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}
