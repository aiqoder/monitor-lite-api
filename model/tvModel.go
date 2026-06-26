package model

import (
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/db"
	"github.com/aiqoder/monitor-lite-api/pkg/common/result"
)

type (
	TvModel struct {
		db         *gorm.DB
		autoOffset int
	}

	Tv struct {
		ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
		Name        string     `json:"title"`        // 原始名称
		DisplayName string     `json:"displayTitle"` // 展示名称
		Url         string     `json:"url" gorm:"uniqueIndex:url_unique"`
		UpdateTime  *time.Time `json:"updateTime" gorm:"type:datetime"`
		Height      uint64     `json:"height"`
		Width       uint64     `json:"width"`
		Speed       uint64     `json:"speed"`
		FailCount   int64      `json:"failCount"`
		Group       string     `json:"group"`  // 节目分组
		Weight      int64      `json:"weight"` // 权重
		ISP         string     `json:"isp"`    // 所属运营商/ 移动、电信等
	}
)

func NewTvModel(cfg *db.Config) TvModel {
	return TvModel{
		db:         cfg.DB,
		autoOffset: 0,
	}
}

func (t *Tv) TableName() string {
	return "tvs"
}

func (t *TvModel) TableList() ([]Tv, error) {
	var tv []Tv

	if err := t.db.Debug().Where("fail_count < ?", 3).Order("weight desc").Order("speed").Order("height desc").Order("update_time desc").Find(&tv).Error; err != nil {
		return tv, err
	}

	return tv, nil
}

// 一键检查专用sql
func (t *TvModel) TableCheckList(typ, extra string) ([]Tv, error) {
	var tv []Tv

	switch typ {
	case "all":
		if err := t.db.Find(&tv).Error; err != nil {
			return tv, err
		}
	case "pix0":
		if err := t.db.Where("width = ? AND height = ?", 0, 0).Find(&tv).Error; err != nil {
			return tv, err
		}
	case "fail":
		if err := t.db.Where("fail_count > ?", 0).Find(&tv).Error; err != nil {
			return tv, err
		}
	case "select":
		if err := t.db.Debug().Where("id in ?", strings.Split(extra, ",")).Find(&tv).Error; err != nil {
			return tv, err
		}
	}

	return tv, nil
}

// 根据名称分组，主要用于更新分组
func (t *TvModel) TableListGroupName() ([]Tv, error) {
	var tv []Tv

	if err := t.db.Where("fail_count < ?", 3).Group("name").Find(&tv).Error; err != nil {
		return tv, err
	}

	return tv, nil
}

// TableListPage 分页查询
func (t *TvModel) TableListPage(m map[string]string) ([]Tv, int64, error) {
	var count int64 = 0
	var tv []Tv
	if err := t.db.Scopes(result.Map2ScopePager(m,
		result.WithFieldZero(map[string]string{
			"width":  "",
			"height": "",
		}),
		result.WithNoLikeKeysKeys("width", "height"))).
		Find(&tv).Error; err != nil {
		return tv, count, err
	}

	t.db.Model(&Tv{}).Scopes(result.Map2ScopeWhere(m, result.WithFieldZero(map[string]string{
		"width":  "",
		"height": "",
	}), result.WithNoLikeKeysKeys("width", "height"))).Count(&count)

	return tv, count, nil
}

func (t *TvModel) TableListByTvName(name string) ([]Tv, error) {
	var tv []Tv

	if err := t.db.Where("fail_count < ? AND name LIKE ?", 5, fmt.Sprintf("%%%s%%", name)).Order("update_time desc").Find(&tv).Error; err != nil {
		return tv, err
	}

	return tv, nil
}

func (t *TvModel) NextTvOffset() ([]Tv, error) {
	var tv []Tv
	err := t.db.Where("fail_count < ?", 3).Order("width").Offset(t.autoOffset).Limit(10).Find(&tv).Error

	t.autoOffset += 1

	if errors.Is(err, gorm.ErrRecordNotFound) {
		t.autoOffset = 0
	}

	if len(tv) == 0 {
		t.autoOffset = 0
	}

	if t.autoOffset < 0 {
		return tv, errors.New("check over")
	}

	return tv, nil
}

func (t *TvModel) Save(tv Tv, fail bool) error {
	existTv := Tv{}
	_ = t.db.Where("url = ?", tv.Url).Find(&existTv)
	if existTv.ID > 0 {
		tv.ID = existTv.ID
		return t.Update(tv, fail)
	}

	if fail {
		return errors.New("检测失败的流")
	}

	ti := time.Now()
	tv.UpdateTime = &ti

	err := t.db.Save(&tv).Error

	if err != nil {
		log.Error(err)
	}
	return err
}

func (t *TvModel) BatchSave(tv []Tv) error {
	return t.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(tv, 10).Error
}

func (t *TvModel) Update(tv Tv, fail bool) error {
	if fail {
		go func() {
			t.db.Table("tvs").Where("id = ?", tv.ID).Update("fail_count", gorm.Expr("fail_count + 1"))
		}()
	}

	ti := time.Now()
	tv.UpdateTime = &ti
	return t.db.Updates(&tv).Error
}

// EmptyGroup 清空分组
func (t *TvModel) EmptyGroup() error {
	return t.db.Table("tvs").Where("1 = ?", 1).Update("display_name", "").Update("group", "").Error
}

func (t *TvModel) Delete(ids []string) {
	if len(ids) == 1 && ids[0] == "-1" {
		_ = t.db.Where("1 = ?", 1).Delete(&Tv{})
	} else {
		for i := 0; i < len(ids); i++ {
			_ = t.db.Where("id = ?", ids[i]).Delete(&Tv{})
		}
	}
}

func (t *TvModel) DeleteAllFail() error {
	return t.db.Where("fail_count > 0").Delete(&Tv{}).Error
}

func (t *TvModel) Updates(tv Tv) error {
	return t.db.Updates(&tv).Error
}

func (t *TvModel) UpdatesByName(tv Tv) error {
	return t.db.Table("tvs").Where("name = ?", tv.Name).Updates(map[string]interface{}{"display_name": tv.DisplayName, "group": tv.Group}).Error
}

func (t *TvModel) Take(id uint64) Tv {
	var tv Tv
	t.db.Where("id = ?", id).Take(&tv)
	return tv
}

func Str2Tv(s string) Tv {
	t := Tv{}
	_ = json.Unmarshal([]byte(s), &t)
	return t
}

// 查询分辨率列表
func (t *TvModel) Pixs() []Tv {
	var tvs []Tv
	err := t.db.Select("width", "height").Order("width").Find(&tvs).Error

	if err != nil {
		log.Error(err)
		return nil
	}

	return tvs
}
