package core

import (
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/protos"
)

// QuestHandler 定义任务处理器契约。
type QuestHandler interface {
	Init(player *playerdomain.Player, quest *playerdomain.Quest)
	GetQuestType() int32
	OnQuestFinish(player *playerdomain.Player, quest *playerdomain.Quest)
	CheckProgress(player *playerdomain.Player, quest *playerdomain.Quest)
	SubscribeEvent()
	HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any)
}

// QuestDirector 定义任务类别切面控制器契约。
type QuestDirector interface {
	OnPlayerLogin(player *playerdomain.Player)
	AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest)
	OnQuestProgressChanged(player *playerdomain.Player, quest *playerdomain.Quest)
	OnQuestProgressFinished(player *playerdomain.Player, quest *playerdomain.Quest)
	TakeProgressRewards(player *playerdomain.Player) []*protos.RewardVo
	GetCategoryType() int32
}

// Resolver 负责 handler 与 director 的查找，避免跨层直接依赖全局函数。
type Resolver interface {
	GetQuestDirector(catalog int32) QuestDirector
	GetQuestHandlerByType(questType int32) QuestHandler
	AcceptQuest(player *playerdomain.Player, questID int32) (*playerdomain.Quest, error)
	ResetQuests(player *playerdomain.Player, category int32)
}

// ResolverAware 表示组件支持注入 Resolver。
type ResolverAware interface {
	SetResolver(resolver Resolver)
}
