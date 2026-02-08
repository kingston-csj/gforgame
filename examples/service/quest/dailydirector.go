package quest

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/protos"
)

// 每日任务类别
type DailyQuestDirector struct {
	baseQuestDirector
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

func (d *DailyQuestDirector) TakeProgressRewards(player *playerdomain.Player) []*protos.RewardVo {
	rewardIndex := player.DailyReset.QuestDailyRewardIndex
	myScore := player.DailyReset.DailyQuestScore
	if myScore > 100 {
		myScore = 100
	}
	canRewardIndex := myScore /20

	commonContainer := config.GetSpecificContainer[*container.CommonContainer]()
	rewardStr := commonContainer.GetStringValue(constants.CommonValueKeyDailyQuestProessReward)
	rewardList := reward.ParseRewardList(rewardStr)
	rewardVos := make([]*protos.RewardVo, 0)
	andReward := reward.NewAndReward()
	for i:=rewardIndex+1;i<=canRewardIndex;i++{
		r := rewardList[i]
		andReward.AddReward(r)
		rewardVos = append(rewardVos, reward.ToRewardVo(r))
	}
	andReward = andReward.Merge()
	andReward.Reward(player, constants.ActionType_DailyQuestProgressReward)
	player.DailyReset.QuestDailyRewardIndex = canRewardIndex

	return rewardVos
}
 
// 获取任务类型
func (d *DailyQuestDirector) GetCategoryType() int32 {
	return  constants.QuestCategoryDaily
}
 







