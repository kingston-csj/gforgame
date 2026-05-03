package bootstrap

import (
	heroservice "io/github/gforgame/examples/service/hero"
	itemservice "io/github/gforgame/examples/service/item"
)

// InitServices 预热服务并完成跨模块注册（reward/consume ops 等）。
func InitServices() {
	heroservice.GetHeroService()
	itemservice.GetItemService()
}
