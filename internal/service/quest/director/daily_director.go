package director

import (
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/config/container"
	"github.com/forfun/gforgame/internal/constants"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/reward"
)

// 每日任务类别
type DailyQuestDirector struct {
	*baseQuestDirector
}

func NewDailyQuestDirector() *DailyQuestDirector {
	d := &DailyQuestDirector{}
	d.baseQuestDirector = NewBaseQuestDirector(d)
	return d
}

/// 实现QuestDirector接口
func (d *DailyQuestDirector) OnPlayerLogin(player *playerdomain.Player) {
	questBox := player.QuestBox
	anyQuests := questBox.SelectUnFinishedQuestsByCategory(constants.QuestCategoryDaily)
	if len(anyQuests) == 0 {
		// 重置任务
		d.resolver.ResetQuests(player, constants.QuestCategoryDaily)
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
	commonContainer := config.GetSpecificContainer[*container.CommonContainer]()
	// 2012_1,2012_1
	maxScoreSum := commonContainer.GetInt32Value(constants.CommonValueKeyDailyQuestScoreSum)
	if myScore > maxScoreSum {
		myScore = maxScoreSum
	}
	
	rewardStr := commonContainer.GetStringValue(constants.CommonValueKeyDailyQuestProessReward)
	rewardList := reward.ParseRewardList(rewardStr)
	canRewardIndex := myScore / (maxScoreSum /int32(len(rewardList)))
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
 







