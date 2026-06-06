package bootstrap

import "github.com/forfun/gforgame/network"

// InitRouteModules 初始化并注册所有路由模块。
func InitRouteModules(router *network.MessageRoute, modules []any) error {
	for _, module := range modules {
		if err := router.RegisterMessageHandlers(module); err != nil {
			return err
		}
	}
	return nil
}
