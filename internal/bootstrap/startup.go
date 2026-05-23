package bootstrap

import (
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"

	dataconfig "github.com/forfun/gforgame/internal/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/system"
)

// 生成顺序必须固定：
// 1) 先生成协议导出与 register_gen.go
// 2) 再基于 register_gen.go 生成 route_dispatch_gen.go
//go:generate go run ../../tools/protocolgen
//go:generate go run ../../tools/routedispatch

// InitMysqlDdl 初始化基础设施（数据库表结构等）。
func InitMysqlDdl() {
	err := mysqldb.Db.AutoMigrate(&playerdomain.Player{})
	if err != nil {
		panic(err)
	}
	err = mysqldb.Db.AutoMigrate(&playerdomain.Friend{})
	if err != nil {
		panic(err)
	}
	err = mysqldb.Db.AutoMigrate(&playerdomain.Scene{})
	if err != nil {
		panic(err)
	}
	err = mysqldb.Db.AutoMigrate(&system.SystemParameterEnt{})
	if err != nil {
		panic(err)
	}
}

// InitConfig 初始化配置数据。
func InitConfig() {
	dataconfig.GetDataManager()
}


// InitBusiness 预热业务数据和计划任务。
func InitBusiness(s *Services) {
	s.Player.LoadPlayerProfile()
	s.Activity.ScheduleAllActivity()
}

// StartSchedulers 启动系统级定时任务。
func StartSchedulers() {
	system.StartSystemTask()
}
