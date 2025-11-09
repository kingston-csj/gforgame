package quest

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)
type DailyQuestDirector struct {
}

func NewDailyQuestDirector() *DailyQuestDirector {
	return &DailyQuestDirector{}
}

/// 实现QuestDirector接口

func (d *DailyQuestDirector) OnPlayerLogin(player *playerdomain.Player) {
	questBox := player.QuestBox
	anyQuests := questBox.SelectUnFinishedQuestsByCategory(constants.QuestCategoryDaily)
	if len(anyQuests) == 0 {
		// 重置任务
		GetQuestService().ResetQuests(player, constants.QuestCategoryDaily)
	}
	quests := questBox.SelectUnFinishedQuestsByCategory(constants.QuestCategoryDaily)
	questVos := make([]*protos.QuestVo, 0, len(quests))
	for _, quest := range quests {
		questVos = append(questVos, quest.ToVo())
	}
	notify := &protos.PushQuestDailyInfo {
		DailyRewardIndex: player.DailyReset.QuestDailyRewardIndex,
		Quests: questVos,
		DailyScore: player.DailyReset.DailyQuestScore,
	}
	io.NotifyPlayer(player, notify)
}

// 玩家完成任务后触发
func (d *DailyQuestDirector) AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest) {
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	player.DailyReset.DailyQuestScore += questData.Score
}
 
// 获取任务类型
func (d *DailyQuestDirector) GetCategoryType(quest *playerdomain.Quest) int {
	return int(constants.QuestCategoryDaily)
}
 







