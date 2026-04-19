package persistence

import (
	"io/github/gforgame/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db *gorm.DB
)

func init() {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	var database, err = gorm.Open(mysql.Open(config.ServerConfig.DbUrl), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}
	Db = database
}
