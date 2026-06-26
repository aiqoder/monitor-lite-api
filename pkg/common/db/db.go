package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

type Config struct {
	DB *gorm.DB
}

type Option func(*Config)

func NewMysql(source string) Option {
	var sqlLogger logger.Interface
	if os.Getenv("MODE") == "dev" {
		sqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		sqlLogger = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		QueryFields: true,
		Logger:      sqlLogger,
	})

	if err != nil {
		panic(err)
	}

	return func(config *Config) {
		config.DB = db
	}
}

type Pager struct {
	Current int
	Size    int
	Total   int64
}

func Pagination(pager *Pager) func(db *gorm.DB) *gorm.DB {
	if pager.Current <= 0 {
		pager.Current = 1
	}

	if pager.Size <= 0 {
		pager.Size = 10
	}

	if pager.Size > 50 {
		pager.Size = 50 // 一次最多查询20条数据
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((pager.Current - 1) * pager.Size).Limit(pager.Size)
	}
}

func NewDB(options ...Option) *Config {
	cfg := &Config{}
	for _, opt := range options {
		opt(cfg)
	}

	return cfg
}
