package bootstrap

import (
	"os"

	mysqldb "io/github/gforgame/examples/infra/persistence"
	protocolexporter "io/github/gforgame/tools/protocol"

	dataconfig "io/github/gforgame/examples/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/service/activity"
	"io/github/gforgame/examples/service/player"
	"io/github/gforgame/examples/system"
	"io/github/gforgame/network"
)

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
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	if env != "dev" {
		return
	}

	generator := protocolexporter.NewCSharpGenerator(
		"examples\\protos",
		"tools\\protocol\\output\\csharp\\",
		"tools\\protocol\\templates\\csharptemplate.tpl",
	)
	if err := generator.Generate(network.GetMsgName2IdMapper()); err != nil {
		panic(err)
	}
	if err := generator.BaseGenerator.GenerateRegisterFromTags("examples\\protos", "examples\\protos\\register_gen.go", nil); err != nil {
		panic(err)
	}
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
