package player

import configcontract "github.com/forfun/gforgame/internal/config/contracts"

// ItemConfigProviders 聚合玩家各背包所需的配置提供器。
// 作为加载期依赖（AfterLoad/Reset 的入参）注入，避免在 Player 实体上持有 infra 状态。
type ItemConfigProviders struct {
	Base  configcontract.ItemConfigProvider
	Rune  configcontract.ItemConfigProvider
	Scene configcontract.ItemConfigProvider
	Equip configcontract.ItemConfigProvider
	Card  configcontract.ItemConfigProvider
}
