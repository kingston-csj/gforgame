package persistence

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/forfun/gforgame/internal/serverconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db     *gorm.DB
	dbOnce sync.Once
)

// InitMysql 初始化数据库连接，必须在 serverconfig.Init() 之后调用。
func InitMysql() {
	dbOnce.Do(func() {
		dbURL, ok := serverconfig.GetExtraString("db.url")
		if !ok || dbURL == "" {
			panic("配置项 db.url 为空，请在 config.yml 中配置")
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

		database, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{
			Logger: gormLogger,
		})
		if err != nil {
			panic(err)
		}
		Db = database
	})
}
