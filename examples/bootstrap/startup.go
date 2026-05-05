package bootstrap

import (
	mysqldb "io/github/gforgame/examples/infra/persistence"

	dataconfig "io/github/gforgame/examples/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/service/activity"
	"io/github/gforgame/examples/service/player"
	"io/github/gforgame/examples/system"
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

// DevOnlyExportProtocols 开发环境导出客户端协议与注册代码。
func DevOnlyExportProtocols() {
	// 已迁移到 go:generate，运行时不再生成文件，避免启动抖动和二次启动副作用。
}

// InitConfig 初始化配置数据。
func InitConfig() {
	dataconfig.GetDataManager()
}

// InitBusiness 预热业务数据和计划任务。
func InitBusiness() {
	player.GetPlayerService().LoadPlayerProfile()
	activity.GetActivityService().ScheduleAllActivity()
}

// StartSchedulers 启动系统级定时任务。
func StartSchedulers() {
	system.StartSystemTask()
}
