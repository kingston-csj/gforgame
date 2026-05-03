package quest

import (
	"sync"
	"time"

	commonerrors "io/github/gforgame/common/errors"
	"io/github/gforgame/common/util"
	"io/github/gforgame/common/util/timeutil"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/protos"
	"io/github/gforgame/examples/reward"
	qcore "io/github/gforgame/examples/service/quest/core"
	questdirector "io/github/gforgame/examples/service/quest/director"
	questhandler "io/github/gforgame/examples/service/quest/handler"
)

// 任务模块
type QuestService struct {
	directors map[int32]qcore.QuestDirector
	handlers  map[int32]qcore.QuestHandler
}

var (
	instance *QuestService
	once     sync.Once
)

func GetQuestService() *QuestService {
	once.Do(func() {
		instance = &QuestService{}

		// 注册所有任务分类
		instance.directors = make(map[int32]qcore.QuestDirector)
		instance.directors[constants.QuestCategoryDaily] = questdirector.NewDailyQuestDirector()
		instance.directors[constants.QuestCategoryMain] = questdirector.NewMainQuestDirector()
		instance.directors[constants.QuestCategoryAchievement] = questdirector.NewAchievementQuestDirector()

		// 注册所有任务类型
		instance.handlers = make(map[int32]qcore.QuestHandler)
		instance.handlers[constants.QuestTypeRecruit] = &questhandler.RecruitQuestHandler{}
		instance.handlers[constants.QuestTime] = &questhandler.TimeQuestHandler{}
		instance.handlers[constants.QuestTypeHeroUpLevel] = &questhandler.HeroLevelUpQuestHandler{}
		instance.handlers[constants.QuestTypeGoldConsume] = &questhandler.GoldConsumeQuestHandler{}
		instance.handlers[constants.QuestTypeDiamondConsume] = &questhandler.DiamondConsumeQuestHandler{}
		instance.handlers[constants.QuestTypeFuben] = &questhandler.FubenLevelQuestHandler{}
		instance.handlers[constants.QuestTypeEquipUpLevel] = &questhandler.EquipUpLevelQuestHandler{}
		instance.handlers[constants.QuestTypeLogin] = &questhandler.LoginQuestHandler{}
		instance.handlers[constants.QuestTypePassGuanka] = &questhandler.MainGuanKaQuestHandler{}

		instance.wireResolver()
		for _, handler := range instance.handlers {
			handler.SubscribeEvent()
		}
	})
	return instance
}

func (s *QuestService) wireResolver() {
	for _, handler := range s.handlers {
		if aware, ok := handler.(qcore.ResolverAware); ok {
			aware.SetResolver(s)
		}
	}
	for _, director := range s.directors {
		if aware, ok := director.(qcore.ResolverAware); ok {
			aware.SetResolver(s)
		}
	}
}

func (s *QuestService) OnPlayerLogin(player *playerdomain.Player) {
	for id, quest := range player.QuestBox.Doing {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		// 删除已完成的任务
		if questData == nil {
			delete(player.QuestBox.Doing, id)
		}
	}
	for catalog := range instance.directors {
		instance.directors[catalog].OnPlayerLogin(player)
	}
}

func (s *QuestService) OnPlayerDailyReset(player *playerdomain.Player) {
	player.QuestBox.ClearQuestsByCategory(constants.QuestCategoryDaily)
	instance.directors[constants.QuestCategoryDaily].OnPlayerLogin(player)
}

func (s *QuestService) ResetQuests(player *playerdomain.Player, catalog int32) {
	questBox := player.QuestBox
	questBox.ClearQuestsByCategory(catalog)
	c := config.GetSpecificContainer[*container.QuestContainer]()
	quests := c.GetRecordsByIndex("Category", catalog)
	for _, quest := range quests {
		s.AcceptQuest(player, quest.Id)
	}
	context.EventBus.Publish(events.PlayerEntityChange, player)
}

func (s *QuestService) AcceptQuest(player *playerdomain.Player, questId int32) (*playerdomain.Quest, error) {
	if player.QuestBox.HasReceivedQuest(questId) {
		return nil, commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	questData := config.QueryById[configdomain.QuestData](questId)
	if questData == nil {
		return nil, commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}

	quest := &playerdomain.Quest{
		Id: questId,
	}
	handler, ok := instance.handlers[questData.Type]
	// 判空（暂时不处理）
	if !ok {
		return nil, commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	handler.Init(player, quest)
	// 有些任务,例如主角升至100,可能接到任务的时已经完成, 需要检查进度
	if quest.IsComplete() {
		handler.CheckProgress(player, quest)
	}

	player.QuestBox.AcceptNewQuest(quest)
	return quest, nil
}

func (s *QuestService) TakeReward(player *playerdomain.Player, questId int32) (*protos.ResQuestTakeReward, *commonerrors.BusinessError) {
	quest := player.QuestBox.GetQuest(questId)
	if quest == nil || quest.Status == constants.QuestStatusRewarded {
		return nil, commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	questData := config.QueryById[configdomain.QuestData](questId)
	rewards := reward.ParseReward(questData.Rewards)
	rewards.Reward(player, constants.ActionType_QuestReward)
	quest.Status = constants.QuestStatusRewarded

	s.GetQuestDirector(questData.Category).AfterTakeReward(player, quest)
	// 主线任务，才要放到完成队列
	if questData.Category == int32(constants.QuestCategoryMain) {
		player.QuestBox.AddFinishedQuest(questId)
	}

	context.EventBus.Publish(events.PlayerEntityChange, player)

	response := &protos.ResQuestTakeReward{
		DailyScore:  int32(player.DailyReset.DailyQuestScore),
		WeeklyScore: int32(player.WeeklyReset.WeeklyQuestScore),
		RewardVos:   reward.ToRewardVos(rewards),
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
		return nil, commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}

	finishedIds := make([]int32, 0)
	for _, quest := range quests {
		finishedIds = append(finishedIds, quest.Id)
		quest.Status = constants.QuestStatusRewarded
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		s.GetQuestDirector(questData.Category).AfterTakeReward(player, quest)
		finishedIds = append(finishedIds, quest.Id)
		// 主线任务，才要放到完成队列
		if questData.Category == int32(constants.QuestCategoryMain) {
			player.QuestBox.AddFinishedQuest(quest.Id)
		}
	}
	rewards.Reward(player, constants.ActionType_QuestRewardAll)
	questBox.ClearQuestsByCategory(catalog)
	context.EventBus.Publish(events.PlayerEntityChange, player)
	response := &protos.ResQuestTakeAllRewards{
		DailyScore:  int32(player.DailyReset.DailyQuestScore),
		WeeklyScore: int32(player.WeeklyReset.WeeklyQuestScore),
		RewardVos:   reward.ToRewardVos(rewards),
	}
	return response, nil
}

func (s *QuestService) GmFinish(player *playerdomain.Player, questId int32) {
	questBox := player.QuestBox
	if questBox.HasReceivedQuest(questId) {
		quest := questBox.Doing[questId]
		questData := quest.Prototype()
		quest.Status = constants.QuestStatusFinished
		quest.Progress = questData.Target
		questDirector := s.GetQuestDirector(questData.Category)
		questDirector.OnQuestProgressChanged(player, quest)
	}
}

func (s *QuestService) TakeProgressReward(player *playerdomain.Player, category int32) *protos.ResQuestTakeProgressReward {
	response := &protos.ResQuestTakeProgressReward{}
	rewardVos := s.GetQuestDirector(category).TakeProgressRewards(player)

	if category == int32(constants.QuestCategoryDaily) {
		response.RewardIndex = player.DailyReset.QuestDailyRewardIndex
	} else if category == int32(constants.QuestCategoryWeekly) {
		response.RewardIndex = player.WeeklyReset.QuestWeeklyRewardIndex
	}

	response.RewardVos = rewardVos
	return response
}

func (s *QuestService) TakeAllRewards(player *playerdomain.Player, catalog int32) (*protos.ResQuestTakeAllRewards, int32) {
	questBox := player.QuestBox
	quests := make([]*playerdomain.Quest, 0)
	questPool := questBox.SelectUnFinishedQuestsByCategory(catalog)
	for _, quest := range questPool {
		if quest.Status != constants.QuestStatusFinished {
			continue
		}
		quests = append(quests, quest)
	}
	if len(quests) == 0 {
		return nil, int32(constants.I18N_QUEST_QUEST_NOT_FINISHED)
	}
	andReward := reward.NewAndReward()
	for _, quest := range quests {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		andReward.AddReward(reward.ParseReward(questData.Rewards))
	}
	andReward = andReward.Merge()
	if err := andReward.Verify(player); err != nil {
		return nil, int32(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	finishedIds := make([]int32, 0)
	score := 0
	for _, quest := range quests {
		quest.Status = constants.QuestStatusRewarded
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		s.GetQuestDirector(questData.Category).AfterTakeReward(player, quest)
		finishedIds = append(finishedIds, quest.Id)
		score += int(questData.Score)
		if questData.Category == int32(constants.QuestCategoryMain) {
			player.QuestBox.AddFinishedQuest(quest.Id)
		}
	}
	andReward.Reward(player, constants.ActionType_QuestRewardAll)

	context.EventBus.Publish(events.PlayerEntityChange, player)
	response := &protos.ResQuestTakeAllRewards{
		DailyScore:  int32(player.DailyReset.DailyQuestScore),
		WeeklyScore: int32(player.WeeklyReset.WeeklyQuestScore),
		RewardVos:   reward.ToRewardVos(andReward),
		Score:       int32(score),
		QuestIds:    finishedIds,
	}
	return response, 0
}

func (s *QuestService) EntrustQuest(player *playerdomain.Player, questId int32, heroId int32) int32 {
	questBox := player.QuestBox
	quest := questBox.GetQuest(questId)
	if quest == nil {
		return int32(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	hero := player.HeroBox.GetHero(heroId)
	if hero == nil {
		return int32(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	if hero.EntrustQuestId != 0 {
		return int32(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	quest.AcceptTime = time.Now().UnixMilli()
	hero.EntrustQuestId = questId
	quest.Status = constants.QuestStatusDoing
	questData := config.QueryById[configdomain.QuestData](questId)

	context.EventBus.Publish(events.PlayerEntityChange, player)

	context.TaskScheduler.Schedule(func() {
		handler, ok := instance.handlers[questData.Type]
		if ok {
			quest.AddProgress(1)
			handler.CheckProgress(player, quest)
		}
	}, timeutil.MILLIS_PER_SECOND*int64(util.Int32Value(questData.Extra)))
	return 0
}

func (s *QuestService) GetQuestDirector(catalog int32) qcore.QuestDirector {
	item, ok := s.directors[catalog]
	if !ok {
		panic("quest director not found")
	}
	return item
}

func (s *QuestService) GetQuestHandlerByType(questType int32) qcore.QuestHandler {
	handler, ok := s.handlers[questType]
	if ok {
		return handler
	}
	return nil
}
