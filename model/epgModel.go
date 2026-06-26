package model

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/gzipx"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"encoding/xml"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slices"
	"strings"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/db"
)

type (
	EpgModel struct {
		db *gorm.DB
	}

	Epg struct {
		ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
		Title       string     `json:"title"`
		Date        string     `json:"date"`
		Start       string     `json:"start"`
		Stop        string     `json:"stop"`
		Desc        string     `json:"desc"`
		Channel     string     `json:"channel"`
		UpdatedTime *time.Time `json:"updated_time"`
	}

	EpgData struct {
		Channel string `xml:"channel,attr"` // 此处的channel是EpgChannel当中的ID
		Start   string `xml:"start,attr"`
		End     string `xml:"stop,attr"`
		Title   string `xml:"title"`
		Desc    string `xml:"desc"`
		Date    string `xml:"date"`
	}

	EpgChannel struct {
		ID          string `xml:"id,attr"`
		DisplayName string `xml:"display-name"`
	}

	EpgRoot struct {
		XMLName   xml.Name     `xml:"tv"`
		Programme []EpgData    `xml:"programme"`
		Channel   []EpgChannel `xml:"channel"`
	}

	Rule struct {
		Group map[string]string `yaml:"group"`
		Name  map[string]string `yaml:"name"`
	}
)

func NewEpgModel(cfg *db.Config) EpgModel {
	return EpgModel{
		db: cfg.DB,
	}
}

func (e *EpgModel) ruleContent() string {
	var s Setting
	e.db.Where("key = ?", "rule").Take(&s)
	return s.Value
}

func (t *Epg) TableName() string {
	return "epg"
}

var epgUrl = []string{
	"https://epg.v1.mk/fy.xml",
	"https://epg.112114.eu.org/pp.xml.gz",
	"https://epg.pw/xmltv/epg_CN.xml.gz",
	"https://epg.pw/xmltv/epg_HK.xml.gz",
	"https://epg.pw/xmltv/epg_TW.xml.gz",
	"http://epg.51zmt.top:8000/difang.xml.gz",
	"http://epg.51zmt.top:8000/cc.xml.gz",
	"https://epg.erw.cc/e.xml.gz",
	"https://epg.pw/xmltv/epg_CN.xml.gz",
}

// UpdateXml2db 将epg数据存入数据库
func (e *EpgModel) UpdateXml2db() {
	toString := e.ruleContent()
	rule := Rule{}
	yaml.Unmarshal([]byte(toString), &rule)

	tvNameArray := maputil.Keys(rule.Name)
	// 储存已经入库的节目
	var cachedChannel []string
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	e.db.Model(&Epg{}).Where("updated_time >= ?", midnight).Group("channel").Pluck("channel", &cachedChannel)

	// 删除早于7天的数据
	err := e.db.Where("updated_time <= ?", datetime.AddDay(time.Now(), -7)).Delete(&Epg{}).Error

	if err != nil {
		log.Error(err)
	}

	for i := 0; i < len(epgUrl); i++ {
		url := epgUrl[i]
		get, err := resty.New().R().Get(url)

		if err != nil {
			log.Error(err)
			continue
		}

		if strings.HasSuffix(url, ".gz") {
			var ed EpgRoot
			xmlData, err := gzipx.Gunzip(get.Body())
			if err != nil {
				log.Error(err)
				continue
			}

			err = xml.Unmarshal(xmlData, &ed)

			if err != nil {
				log.Error(err)
				continue
			}

			// 用户临时存储id和最终的DisplayName，用户加快构建速度
			idEpgMap := map[string]string{}

			getDisplayName := func(id string) string { // ID 对应 programme 当中的 channel字段
				displayName, exist := idEpgMap[id]

				if exist {
					return displayName
				}

				// 特殊频道补偿措施
				epgName := map[string]string{
					"CCTV1": "CCTV-1综合",
				}

				for _, ch := range ed.Channel {
					if ch.ID == id {
						displayName = ch.DisplayName

						// 容易出错的频道
						s, ok := epgName[displayName]

						if ok {
							return s
						}

						// 没有找到从规则当中补偿
						if !slices.Contains(tvNameArray, displayName) {
							for key, value := range rule.Name {
								if strings.Contains(strings.ToLower(value), strings.ToLower(displayName)) {
									displayName = key
									break
								}
							}
						}
						idEpgMap[id] = displayName
						return displayName
					}
				}

				return ""
			}

			programme := ed.Programme

			with := slice.GroupWith(programme, func(item EpgData) string {
				displayName := getDisplayName(item.Channel)
				return displayName
			})

			for i := 0; i < len(tvNameArray); i++ {
				tvName := tvNameArray[i]
				epgSlice, ok := with[tvName]
				if ok {
					// 已存在的数据，如果该频道已经缓存，则跳过不如库
					if slices.Contains(cachedChannel, tvName) {
						continue
					}
					var insertEpg []Epg
					for i := 0; i < len(epgSlice); i++ {
						epg := epgSlice[i]
						//tempStoreArray = append(tempStoreArray, epg.Title)
						layout := "20060102150405  +0800"
						startTime, err := time.Parse(layout, epg.Start)
						if err != nil {
							continue
						}
						stopTime, err := time.Parse(layout, epg.End)
						if err != nil {
							continue
						}

						now := time.Now()

						date := startTime.Format(time.DateOnly)

						insertEpg = append(insertEpg, Epg{
							Title:       epg.Title,
							Date:        date,
							Start:       startTime.Format(time.TimeOnly),
							Stop:        stopTime.Format(time.TimeOnly),
							Channel:     tvName,
							Desc:        epg.Desc,
							UpdatedTime: &now,
						})
					}

					if len(insertEpg) > 0 {
						batchErr := e.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(insertEpg, len(insertEpg)).Error
						if batchErr != nil {
							log.Error(batchErr)
						}
					}
				}
			}
		}
	}
}

func (e *EpgModel) List(channel string) []Epg {
	var epgs []Epg
	err := e.db.Where("channel = ? AND updated_time >= ?", channel, datetime.AddDay(time.Now(), -7)).Find(&epgs).Error

	if err != nil {
		log.Error(err)
		return nil
	}

	return epgs
}

func (e *EpgModel) ChannelEpgByDate(channel, date string) []Epg {
	var epgs []Epg
	err := e.db.Where("channel = ? AND date = ?", channel, date).Find(&epgs).Error

	if err != nil {
		log.Error(err)
		return nil
	}

	return epgs
}
