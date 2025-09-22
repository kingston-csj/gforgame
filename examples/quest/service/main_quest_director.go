package quest

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
)
type MainQuestDirector struct {
}

func NewMainQuestDirector() *MainQuestDirector {
	return &MainQuestDirector{}
}

/// 实现QuestDirector接口

func (d *MainQuestDirector) OnPlayerLogin(player *playerdomain.Player) {
	firstMainQuestId := config.GetSpecificContainer[configdomain.QuestData, container.QuestContainer]("quest").FirstMainQuestId
	if !player.QuestBox.HasReceivedQuest(firstMainQuestId) {
		// player.QuestBox.AddQuest(firstMainQuestId)
	}
}

// 玩家完成任务后触发
func (d *MainQuestDirector) AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest) {
	// 1. 从数据库加载玩家任务
	// 2. 转换为vo
	// 3. 下发给客户端
}

// 任务进度变更触发
func (d *MainQuestDirector) OnQuestProgressChanged(player *playerdomain.Player, quest *playerdomain.Quest) {
	// 1. 从数据库加载玩家任务
	// 2. 转换为vo
	// 3. 下发给客户端
}

// 获取任务类型
func (d *MainQuestDirector) GetCategoryType(quest *playerdomain.Quest) int {
	return int(constants.QuestCategoryMain)
}
 







