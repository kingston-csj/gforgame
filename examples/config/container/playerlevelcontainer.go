package container

import (
	"io/github/gforgame/data"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/attribute"
)

type PlayerLevelContainer struct {
	*data.Container[int32, configdomain.PlayerLevelData]
	levelDataMap map[int32]map[int32]*configdomain.PlayerLevelData
	MaxLevel int32
}

func (c *PlayerLevelContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.PlayerLevelData]()
	}
	c.Container.Init()
	c.levelDataMap = make(map[int32]map[int32]*configdomain.PlayerLevelData)
}

// AfterLoad 数据加载后的处理
func (c *PlayerLevelContainer) AfterLoad() {
	c.levelDataMap = make(map[int32]map[int32]*configdomain.PlayerLevelData)
	for _, record := range c.GetAllRecords() {
		modelMap, exists := c.levelDataMap[record.Id]
		if !exists {
			modelMap = make(map[int32]*configdomain.PlayerLevelData)
			c.levelDataMap[record.Id] = modelMap
		}
		modelMap[record.Level] = record
		if record.Level > c.MaxLevel {
			c.MaxLevel = record.Level
		}
	}
}

// GetLevelData 获取指定玩家和等级的数据
func (c *PlayerLevelContainer) GetLevelData(playerId int32, level int32) *configdomain.PlayerLevelData {
	if modelMap, exists := c.levelDataMap[playerId]; exists {
		return modelMap[level]
	}
	return nil
}

// GetHeroLevelAttrs 获取玩家等级属性
func (c *PlayerLevelContainer) GetPlayerLevelAttrs(heroId int32, level int32) []attribute.Attribute {
	data := c.GetLevelData(heroId, level)
	if data == nil {
		return nil
	}
	return data.GetPlayerLevelAttrs()
}
