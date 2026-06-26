package model

import (
	"gorm.io/gorm"
	"github.com/aiqoder/monitor-lite-api/pkg/common/db"
	"github.com/aiqoder/monitor-lite-api/pkg/common/result"
)

//----------------------------------------------------------------

type (
	SelfoutModel struct {
		db *gorm.DB
	}
	Selfout struct {
		Id     int64  `json:"id" gorm:"primaryKey;autoIncrement"`
		Width  int64  `json:"width"`
		Height int64  `json:"height"`
		Speed  int64  `json:"speed"`
		Key    string `json:"key"`
	}
)

func NewSelfoutModel(cfg *db.Config) SelfoutModel {
	return SelfoutModel{db: cfg.DB}
}

func (s *Selfout) TableName() string {
	return "self_out"
}

func (s *SelfoutModel) FirstOne(id int64) (Selfout, error) {
	selfout := Selfout{}
	err := s.db.First(&selfout, "id = ?", id).Error
	return selfout, err
}

func (s *SelfoutModel) FirstOneByKey(key string) (Selfout, error) {
	selfout := Selfout{}
	err := s.db.First(&selfout, "key = ?", key).Error
	return selfout, err
}

func (s *SelfoutModel) List() ([]Selfout, error) {
	selfouts := new([]Selfout)
	err := s.db.Find(&selfouts).Error
	return *selfouts, err
}

// "current": strconv.FormatInt(req.Current, 10),
// "size": strconv.FormatInt(req.Size, 10),
// "id": strconv.FormatInt(req.Id, 10),
// "width": strconv.FormatInt(req.Width, 10),
// "height": strconv.FormatInt(req.Height, 10),
// "speed": strconv.FormatInt(req.Speed, 10),
// "key": req.Key,
func (s *SelfoutModel) ListPage(m map[string]string) (selfout []Selfout, count int64, err error) {
	zero := result.WithFieldZero(map[string]string{
		"width":  "0",
		"height": "0",
		"speed":  "0",
	})
	if err := s.db.Scopes(result.Map2ScopePager(m, zero)).Find(&selfout).Error; err != nil {
		return selfout, count, err
	}

	s.db.Scopes(result.Map2ScopeWhere(m)).Count(&count)
	return selfout, count, nil
}

func (s *SelfoutModel) Save(selfout *Selfout) error {
	err := s.db.Save(&selfout).Error
	return err
}

func (s *SelfoutModel) Del(selfout *Selfout) error {
	return s.db.Model(&Selfout{}).Delete(Selfout{Id: selfout.Id}).Error
}
