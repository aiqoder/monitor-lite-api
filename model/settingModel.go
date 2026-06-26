package model

import (
	"gorm.io/gorm"
	"github.com/aiqoder/monitor-lite-api/pkg/common/db"
)

type (
	SettingModel struct {
		db *gorm.DB
	}

	Setting struct {
		ID    uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
		Key   string `json:"key" gorm:"uniqueIndex:key_unique"` // 原始名称
		Value string `json:"value"`                             // 值
		Type  string `json:"type"`                              // bool、number、string 三种类型
	}
)

func NewSettingModel(cfg *db.Config) SettingModel {
	// 初始化账户
	cfg.DB.Create(&Setting{Key: "username", Value: "admin", Type: "string"})
	// 初始化密码
	cfg.DB.Create(&Setting{Key: "password", Value: "admin123", Type: "string"})
	// 显示未知分组
	cfg.DB.Create(&Setting{Key: "unknownGroup", Value: "0", Type: "bool"})
	// DIYP链接
	cfg.DB.Create(&Setting{Key: "plusKey", Value: "diyp", Type: "string"})
	// 自动删除
	cfg.DB.Create(&Setting{Key: "autoDelCount", Value: "3", Type: "number"})
	// 自动抓取EPG
	cfg.DB.Create(&Setting{Key: "autoEpg", Value: "0", Type: "number"})
	// 自动检测
	cfg.DB.Create(&Setting{Key: "autoCheck", Value: "0", Type: "number"})
	// 自动搜源
	cfg.DB.Create(&Setting{Key: "autoSearch", Value: "0", Type: "string"})
	// 订阅自动读取时间(每天)
	cfg.DB.Create(&Setting{Key: "subscriberTime", Value: "00:00", Type: "string"})
	// 黑名单
	blankListDefValue := "liveplay.myqcloud.com#103.45.68.47#21dtv.com#51daao.com#51romzj.fun#51romzj.shop#51romzj.site#51romzj.xyz#51zb.fun#51zb.work#733240.com#87.98.184.134#8vxbkar.com#8vxcuve.com#8vxdeus.com#8vxevou.com#8vxkewe.com#adultiptv.net#amazzin.pw#awvvvvw.live#bo7521521.com#bybzj.com#bytebwq.com#caomin5168.com#cdn2020.com#cdnedge.live#ckzy1com.com#cntv.sbs#dadi-bo.com#ddyunbo.com#dfkj.live#didivod.com#dlyilian.com#douyu#embedplayer.net#epg.pm#epg.pw#feimanzb.com#goulexizuo.com#hellotvvod.com#huishenghuo888888.com#iaxaa.com#imgscloud.com#iptvxxx#itv.xmtv9527.com#iwant-sex.com#jiufanrj.com#kakoaction.com#kinoprofi.vip#kpkuang.xyz#lbbf9.com#maa1804.com#maa1809.com#mangguo-youku.com#mimivodplay.com#myalicdn#mycamtv.com#mycamtv.net#naibago.com#ottclub.xyz#oywine.com#pelicanhosting.xyz#play92332.com#qcloudcdn#qrneryt.com#quweikm.com#quweikm.com#redtraffic.xyz#rhsj520.com#rongliren.com#saejeuj.com#siwazywcdn2.com#slbfsl.com#stz8.com#suklakakl.site#sytes.net#ttbfp2.com#ttbfp4.com#turner.com#vodyutu.com#whqhyg.com#wzj9.com#yishihui.com#yongaomy.com#zhibotv.cf"
	cfg.DB.Create(&Setting{Key: "blackList", Value: blankListDefValue, Type: "string"})
	cfg.DB.Create(&Setting{Key: "aiBaseUrl", Value: "https://api.openai.com/v1", Type: "string"})
	cfg.DB.Create(&Setting{Key: "aiApiKey", Value: "", Type: "string"})
	cfg.DB.Create(&Setting{Key: "aiModel", Value: "gpt-4o-mini", Type: "string"})
	return SettingModel{
		db: cfg.DB,
	}
}

func (t *SettingModel) RuleContent() string {
	return t.Value("rule")
}

func (t *SettingModel) SaveRule(content string) error {
	var s Setting
	t.db.Where("key = ?", "rule").Take(&s)
	if s.ID == 0 {
		return t.db.Create(&Setting{Key: "rule", Value: content, Type: "yaml"}).Error
	}
	s.Value = content
	return t.db.Save(&s).Error
}

func (t *SettingModel) EnsureRule(content string) {
	if t.RuleContent() != "" {
		return
	}
	_ = t.SaveRule(content)
}

// TableList 查询所有配置
func (t *SettingModel) TableList() ([]Setting, error) {
	var setting []Setting

	if err := t.db.Find(&setting).Error; err != nil {
		return setting, err
	}

	return setting, nil
}

func (t *SettingModel) Update(setting Setting) error {
	return t.db.Updates(setting).Error
}

func (t *SettingModel) UpdateValueByKey(key, value string) error {
	var s Setting
	if err := t.db.Where("key = ?", key).First(&s).Error; err != nil {
		return err
	}
	s.Value = value
	return t.db.Save(&s).Error
}

func (t *SettingModel) Value(key string) string {
	s := Setting{}
	t.db.Where("key = ?", key).Take(&s)
	return s.Value
}
