package model

import (
	"gorm.io/gorm"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/db"
)

type (
	SubscriberModel struct {
		db *gorm.DB
	}

	Subscriber struct {
		ID        uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
		Name      string     `json:"title"` // 原始名称
		Url       string     `json:"url" gorm:"uniqueIndex:url_unique"`
		Count     uint64     `json:"count"`      // 抓取数量
		CheckTime *time.Time `json:"check_time"` // 最后一次抓取时间
	}
)

func NewSubscriberModel(cfg *db.Config) SubscriberModel {
	return SubscriberModel{
		db: cfg.DB,
	}
}

// 查询所有
func (s *SubscriberModel) List() ([]Subscriber, error) {
	var subs []Subscriber
	err := s.db.Find(&subs).Error
	return subs, err
}

func (s *SubscriberModel) Update(subscriber Subscriber) error {
	if subscriber.ID == 0 {
		return s.db.Create(&subscriber).Error
	}
	return s.db.Debug().Updates(&subscriber).Error
}

func (s *SubscriberModel) Delete(subscriber Subscriber) error {
	return s.db.Delete(subscriber).Error
}

func (s *SubscriberModel) GetById(id uint64) (Subscriber, error) {
	sub := Subscriber{}
	err := s.db.Where("id = ?", id).Take(&sub).Error
	return sub, err
}
