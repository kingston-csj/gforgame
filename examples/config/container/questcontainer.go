package container

import (
	"io/github/gforgame/data"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
)

type QuestContainer struct {
	*data.Container[int32, configdomain.QuestData]
	// 主线第一条任务
	FirstMainQuestId int32
}

func (c *QuestContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.QuestData]()
	}
	c.Container.Init()
}

// AfterLoad 数据加载后的处理
func (c *QuestContainer) AfterLoad() {
	for _, record := range c.Values() {
		if record.Category == int32(constants.QuestCategoryMain) {
			if c.FirstMainQuestId == 0 {
				c.FirstMainQuestId = record.Id
			}
		}
	}
}

