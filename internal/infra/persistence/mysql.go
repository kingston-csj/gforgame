package persistence

import (
	"log"
	"os"
	"time"

	"github.com/forfun/gforgame/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db *gorm.DB
)

func init() {
	dbURL, ok := config.GetExtraString("db.url")
	if !ok || dbURL == "" {
		panic("配置项 db.url 为空，请在 config-game.yml 中配置")
	}
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	var database, err = gorm.Open(mysql.Open(dbURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}
	Db = database
}
