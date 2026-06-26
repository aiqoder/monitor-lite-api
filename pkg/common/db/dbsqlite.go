package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func NewSqlite(source string) Option {
	var sqlLogger logger.Interface
	if os.Getenv("MODE") == "dev" {
		sqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		sqlLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(sqlite.Open(source), &gorm.Config{
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
