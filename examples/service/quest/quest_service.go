package quest

import (
	"sync"

	"io/github/gforgame/common"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/reward"
	handler "io/github/gforgame/examples/service/quest/handler"
	"io/github/gforgame/protos"
)

type QuestService struct {
	directors map[int32]QuestDirector 
	handlers map[int32]handler.QuestHandler
}

var (
	instance *QuestService
	once     sync.Once
	
)

func GetQuestService() *QuestService {
	once.Do(func() {
		instance = &QuestService{}

		// 注册所有任务分类
		instance.directors = make(map[int32]QuestDirector)
		// 注册所有任务类型
		instance.handlers = make(map[int32]handler.QuestHandler)

		instance.directors[constants.QuestCategoryDaily] = NewDailyQuestDirector()
		instance.directors[constants.QuestCategoryMain] = NewMainQuestDirector()
		instance.directors[constants.QuestCategoryAchievement] = NewAchievementQuestDirector()

		// 注册所有任务类型
		instance.handlers[constants.QuestTypeRecruit] = &handler.RecruitQuestHandler{}
		for _,handler := range instance.handlers {
			handler.SubscribeEvent()
		}
	})
	return instance
}

func (s *QuestService) OnPlayerLogin(player *playerdomain.Player) {
	for catalog := range instance.directors {
		instance.directors[catalog].OnPlayerLogin(player)
	}
}

func (s *QuestService) ResetQuests(player *playerdomain.Player, catalog int32) {
	questBox := player.QuestBox
	questBox.ClearQuestsByCategory(catalog)
	c := config.GetSpecificContainer[ container.QuestContainer]("quest")
	quests := c.GetRecordsBy("Category", catalog)
	for _, quest := range quests {
		s.AcceptQuest(player, quest.Id)
	}
	context.EventBus.Publish(events.PlayerEntityChange, player)
}

func (s *QuestService) AcceptQuest(player *playerdomain.Player, questId int32) (*playerdomain.Quest, error) {
	if player.QuestBox.HasReceivedQuest(questId) {
		return nil, common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	questData := config.QueryById[configdomain.QuestData](questId)
	if questData == nil {
		return nil, common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
   
	quest := &playerdomain.Quest{
		Id:   questId,
	}
	handler,ok := instance.handlers[questData.Type] 
	// 判空（暂时不处理）
	if ok {
        handler.Init(player, quest)
		 // 有些任务,例如主角升至100级,可能接到任务的时候,已经完成了,需要做个检查
		if quest.IsComplete() {
			handler.CheckProgress(player, quest)
		}
    }

	player.QuestBox.AcceptNewQuest(quest)
	return quest, nil
}

func (s *QuestService) TakeReward(player *playerdomain.Player, questId int32) (*protos.ResQuestTakeReward, *common.BusinessRequestException) {
	quest := player.QuestBox.GetQuest(questId)
	if quest == nil || quest.Status == constants.QuestStatusRewarded{
		return nil, common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	questData := config.QueryById[configdomain.QuestData](questId)
	rewards := reward.ParseReward(questData.Rewards)
	rewards.Reward(player, constants.ActionType_QuestReward)
	quest.Status = constants.QuestStatusRewarded

	GetQuestDirector(questData.Category).AfterTakeReward(player, quest)
	// 主线任务，才要放到完成列表
	if questData.Category == int32(constants.QuestCategoryMain) {
		player.QuestBox.AddFinishedQuest(questId)
	}

	context.EventBus.Publish(events.PlayerEntityChange, player)

	response := &protos.ResQuestTakeReward{
		DailyScore: int32(player.DailyReset.DailyQuestScore),
		WeeklyScore: int32(player.WeeklyReset.WeeklyQuestScore),
		RewardVos: reward.ToRewardVos(rewards),
	}
	return response, nil
}

func (s *QuestService) TakeAllReward(player *playerdomain.Player, catalog int32) (*protos.ResQuestTakeAllRewards, error) {
	questBox := player.QuestBox
	questPool := questBox.SelectUnFinishedQuestsByCategory(catalog)
	quests := make([]*playerdomain.Quest, 0)
	for _, quest := range questPool {
		if quest.Status != constants.QuestStatusFinished {
			continue
		}
		quests = append(quests, quest)
	}
	rewards := reward.NewAndReward()
	for _, quest := range quests {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		rewards.AddReward(reward.ParseReward(questData.Rewards))
	}
	rewards.Merge()
	if err := rewards.Verify(player); err != nil {
		return nil, common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}

	finishedIds := make([]int32, 0)
	for _, quest := range quests {
		finishedIds = append(finishedIds, quest.Id)
		quest.Status = constants.QuestStatusRewarded
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		GetQuestDirector(questData.Category).AfterTakeReward(player, quest)
		finishedIds = append(finishedIds, quest.Id)
		// 主线任务，才要放到完成列表
		if questData.Category == int32(constants.QuestCategoryMain) {
			player.QuestBox.AddFinishedQuest(quest.Id)
		}
	}
	rewards.Reward(player, constants.ActionType_QuestRewardAll)
	questBox.ClearQuestsByCategory(catalog)
	context.EventBus.Publish(events.PlayerEntityChange, player)
	response := &protos.ResQuestTakeAllRewards{
		DailyScore: int32(player.DailyReset.DailyQuestScore),
		WeeklyScore: int32(player.WeeklyReset.WeeklyQuestScore),
		RewardVos: reward.ToRewardVos(rewards),
	}
	return response, nil
}

func (s *QuestService) GmFinish(player *playerdomain.Player, questId int32) {
	questBox := player.QuestBox
	if questBox.HasReceivedQuest(questId) {
		quest := questBox.Doing[questId]
		quest.Status = constants.QuestStatusFinished
		quest.Progress = quest.Target
		questDirector :=  GetQuestDirector(quest.Prototype().Category)
		questDirector.OnQuestProgressChanged(player, quest)
	}
}

func (s *QuestService) TakeProgressReward(player *playerdomain.Player, category int32) *protos.ResQuestTakeProgressReward {
	response := &protos.ResQuestTakeProgressReward{}
    rewardVos := GetQuestDirector(category).TakeProgressRewards(player)

	if category == int32(constants.QuestCategoryDaily) {
		response.RewardIndex =  player.DailyReset.QuestDailyRewardIndex
	} else if category == int32(constants.QuestCategoryWeekly) {
		response.RewardIndex =  player.WeeklyReset.QuestWeeklyRewardIndex
	}

    response.RewardVos = rewardVos
    return response
}

func GetQuestDirector(catalog int32) QuestDirector {
	item,ok := instance.directors[catalog] 
    if !ok {
        panic("quest director not found")
    }
    return item
}
