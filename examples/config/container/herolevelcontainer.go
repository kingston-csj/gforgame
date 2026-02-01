package container

import (
	"io/github/gforgame/data"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/attribute"
)

type HeroLevelContainer struct {
	*data.Container[int32, configdomain.HeroLevelData]
	levelDataMap map[int32]map[int32]*configdomain.HeroLevelData
	MaxLevel int32
}

func (c *HeroLevelContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.HeroLevelData]()
	}
	c.Container.Init()
	c.levelDataMap = make(map[int32]map[int32]*configdomain.HeroLevelData)
}

// AfterLoad 数据加载后的处理
func (c *HeroLevelContainer) AfterLoad() {
	c.levelDataMap = make(map[int32]map[int32]*configdomain.HeroLevelData)
	for _, record := range c.GetAllRecords() {
		modelMap, exists := c.levelDataMap[record.Id]
		if !exists {
			modelMap = make(map[int32]*configdomain.HeroLevelData)
			c.levelDataMap[record.Id] = modelMap
		}
		modelMap[record.Level] = record
		if record.Level > c.MaxLevel {
			c.MaxLevel = record.Level
		}
	}
}

// GetLevelData 获取指定英雄和等级的数据
func (c *HeroLevelContainer) GetLevelData(heroId int32, level int32) *configdomain.HeroLevelData {
	if modelMap, exists := c.levelDataMap[heroId]; exists {
		return modelMap[level]
	}
	return nil
}

// GetHeroLevelAttrs 获取英雄等级属性
func (c *HeroLevelContainer) GetHeroLevelAttrs(heroId int32, level int32) []attribute.Attribute {
	data := c.GetLevelData(heroId, level)
	if data == nil {
		return nil
	}
	return data.GetHeroLevelAttrs()
}
