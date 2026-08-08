package itemconfig

import (
	"github.com/forfun/gforgame/internal/config"
	configcontract "github.com/forfun/gforgame/internal/config/contracts"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
)

// BaseItemConfigProvider 提供普通道具配置查询。
type BaseItemConfigProvider struct{}

func NewBaseItemConfigProvider() configcontract.ItemConfigProvider {
	return &BaseItemConfigProvider{}
}

func (p *BaseItemConfigProvider) GetConfig(itemID int32) configcontract.ItemConfig {
	return config.QueryById[configdomain.PropData](itemID)
}

// RuneConfigProvider 提供符文配置查询。
type RuneConfigProvider struct{}

func NewRuneConfigProvider() configcontract.ItemConfigProvider {
	return &RuneConfigProvider{}
}

func (p *RuneConfigProvider) GetConfig(itemID int32) configcontract.ItemConfig {
	return config.QueryById[configdomain.RuneData](itemID)
}

// SceneItemConfigProvider 提供场景道具配置查询。
type SceneItemConfigProvider struct{}

func NewSceneItemConfigProvider() configcontract.ItemConfigProvider {
	return &SceneItemConfigProvider{}
}

func (p *SceneItemConfigProvider) GetConfig(itemID int32) configcontract.ItemConfig {
	return config.QueryById[configdomain.ScenePropData](itemID)
}
