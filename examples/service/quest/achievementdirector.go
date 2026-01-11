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
type AchievementQuestDirector struct {
	baseQuestDirector
}

func NewAchievementQuestDirector() *AchievementQuestDirector {
	return &AchievementQuestDirector{}
}

/// 实现QuestDirector接口
func (d *AchievementQuestDirector) OnPlayerLogin(player *playerdomain.Player) {
	questBox := player.QuestBox
	questContainer := config.QueryContainer[configdomain.QuestData, *container.QuestContainer]()
	for group, quests := range questContainer.AchievementsByGroup {
		// 初始化每个分组的第一条任务
		if !questBox.HasAppointedTypeAchievement(d.GetCategoryType(), group) {
			GetQuestService().AcceptQuest(player, quests[0].Id)
		}
	}

	achievementVos := make([]*protos.QuestVo, 0)
	for _, achievement := range questBox.SelectUnFinishedQuestsByCategory(d.GetCategoryType()) {
		achievementData := config.QueryById[configdomain.QuestData](achievement.Id)
		// 如果是最后一条，或者未完成状态，则添加到列表
		if (achievementData.Next == 0 || achievement.Status != constants.QuestStatusRewarded) {
			achievementVos = append(achievementVos, achievement.ToVo() )
		}
	}

	notify := &protos.PushAchievementInfo {
		Score: player.ExtendBox.AchievementScore,
		AchievementVos: achievementVos,
	}
	io.NotifyPlayer(player, notify)
}

// 任务完成执行切面
func (d *AchievementQuestDirector) OnQuestProgressFinished(player *playerdomain.Player, quest *playerdomain.Quest) {
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	// 下一个任务，自动继承当前进度
	if (questData.Next > 0) {
		nextQuest, _ := GetQuestService().AcceptQuest(player, questData.Next)
		nextQuest.Progress = quest.Progress
		if (nextQuest.Progress >= nextQuest.Target) {
			nextQuest.Status = constants.QuestStatusFinished
		}
		refresh := &protos.PushQuestRefreshVo{}
		refresh.Quest = nextQuest.ToVo()
		io.NotifyPlayer(player, refresh)
	}
}

// 玩家完成任务后触发
func (d *AchievementQuestDirector) AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest) {
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	player.ExtendBox.AchievementScore += questData.Score
}

func (d *AchievementQuestDirector) TakeProgressRewards(player *playerdomain.Player) []*protos.RewardVo {
	myScore := player.ExtendBox.AchievementScore
	canRewardTimes := myScore /100
	if canRewardTimes < 1 {
		return make([]*protos.RewardVo, 0)
	}

	commonContainer := config.QueryContainer[configdomain.CommonData, *container.CommonContainer]()
	rewardStr := commonContainer.GetStringValue(constants.CommonValueKeyAchievementQuestProessReward)
	perReward := reward.ParseReward(rewardStr)
	rewardVos := make([]*protos.RewardVo, 0)
	andReward := reward.MultiplyAndReward(perReward, float64(canRewardTimes))
	andReward.Reward(player, constants.ActionType_AchievementQuestProgressReward)

	return rewardVos
}
 
// 获取任务类型
func (d *AchievementQuestDirector) GetCategoryType() int32 {
	return  constants.QuestCategoryAchievement
}
 







