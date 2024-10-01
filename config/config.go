package config

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DbUrl     string
	ServerUrl string
	//后端管理地址
	HttpUrl string
	// pprof性能监测地址
	PprofAddr string
}

//go:embed default.yml
var configFS embed.FS

var (
	ServerConfig Config
)

// 多环境配置读取规则：
// 加载顺序: default -> 根据开发环境或者部署环境读取相应子目录config/config-env文件
// 先加载的配置会被后加载的同名配置所替换！！！！
// 1.优先读default.yml文件，应用程序内部配置，项目打包成二进制可执行文件也会嵌入该配置
// 2.在开发环境，若环境变量为xxx，则读取config-xxx.yml(环境变量默认为dev)
// 3.在部署环境（打包成二进制可执行文件），若环境变量为yyy，则读取config-yyy.yml(环境变量默认为dev)
func init() {
	// 创建 Viper 实例
	v := viper.New()
	v.SetConfigType("yml")
	// 打包后的二进制文件也要
	f, err := configFS.Open("default.yml")
	if err != nil {
		panic(fmt.Errorf("failed to open config file: %w", err))
	}
	defer f.Close()
	// 从 io.Reader 读取配置
	if err := v.ReadConfig(f); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	// 允许 Viper 读取环境变量
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// 获取环境变量，确定要加载的配置文件
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	v.SetConfigName("config-" + env)
	// ./config路径在开发及部署环境均适用
	v.AddConfigPath("./config")
	// 再次读取配置文件，这次是根据环境变量，使用合并配置的方法确保旧配置被替换
	if err := v.MergeInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	ServerConfig = Config{
		DbUrl:     v.GetString("db.url"),
		ServerUrl: v.GetString("server.addr"),
		HttpUrl:   v.GetString("server.httpAddr"),
		PprofAddr: v.GetString("server.pprofAddr"),
	}
	fmt.Println("dbUrl", ServerConfig.DbUrl)
	fmt.Println("serverAddr", ServerConfig.ServerUrl)
}
