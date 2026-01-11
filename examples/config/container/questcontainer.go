package container

import (
	"io/github/gforgame/data"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	"sort"
)

type QuestContainer struct {
	*data.Container[int32, configdomain.QuestData]
	// 主线第一条任务
	FirstMainQuestId int32
	// 成就任务按分类分组
	AchievementsByGroup map[int32][]*configdomain.QuestData
}

func (c *QuestContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.QuestData]()
		c.AchievementsByGroup = make(map[int32][]*configdomain.QuestData)
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
		if record.Category == int32(constants.QuestCategoryAchievement) {
			if _, ok := c.AchievementsByGroup[record.Type]; !ok {
				c.AchievementsByGroup[record.Type] = make([]*configdomain.QuestData, 0)
			}
			c.AchievementsByGroup[record.Type] = append(c.AchievementsByGroup[record.Type], record)
		}
	}

	for _, entry := range c.AchievementsByGroup {
		achievements := entry
		sort.Slice(achievements, func(i, j int) bool {
			return achievements[i].Id < achievements[j].Id
		})
		// 设置成就的前后关系
		for i := 0; i < len(achievements); i++ {
			current := achievements[i]
			if i > 0 {
				previous := achievements[i-1]
				current.PreviousId = previous.Id
			} else {
				current.PreviousId = 0; // 第一个成就没有前置
			}
			if (i < len(achievements) - 1) {
				next := achievements[i + 1];
				current.Next = next.Id;
			} else {
				current.Next = 0; // 最后一个成就没有后置
			}
		}
	}
}

