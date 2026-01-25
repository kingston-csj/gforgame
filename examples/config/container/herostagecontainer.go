package container

import (
	"io/github/gforgame/data"
	configdomain "io/github/gforgame/examples/domain/config"
)

type HeroStageContainer struct {
	*data.Container[int32, configdomain.HeroStageData]
	stageMap map[int32]*configdomain.HeroStageData
}

func (c *HeroStageContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.HeroStageData]()
	}
	c.Container.Init()
	c.stageMap = make(map[int32]*configdomain.HeroStageData)
}

// AfterLoad 数据加载后的处理
func (c *HeroStageContainer) AfterLoad() {
	c.stageMap = make(map[int32]*configdomain.HeroStageData)
	// 使用新增的GetAllRecords方法来访问所有记录
	for _, record := range c.GetAllRecords() {
		c.stageMap[record.Stage] = record
	}
}

// GetRecordByStage 根据关卡获取数据
func (c *HeroStageContainer) GetRecordByStage(stage int32) *configdomain.HeroStageData {
	return c.stageMap[stage]
}
