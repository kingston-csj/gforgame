package quest

import (
	"io/github/gforgame/examples/config"
	constants "io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
	"io/github/gforgame/util"
)
var (
	// 绑定各种类型的任务处理器
	handlers = map[int32]QuestHandler{}
)

func init() {
	handlers[constants.QuestTypeHeroSum] = &HeroSumQuestHandler{}
	// 初始化，订阅事件
	for _, h := range handlers {
		h.SubscribeEvent()
	}
}

type QuestHandler interface {
	// 任务初始化，主要是进度相关
	Init(player *playerdomain.Player,quest *playerdomain.Quest)

	GetQuestType() int32
	// 任务完成时调用
	OnQuestFinish(player *playerdomain.Player,quest *playerdomain.Quest)

	// 检查任务进度
	CheckProgress(player *playerdomain.Player,quest *playerdomain.Quest) 

	// 订阅事件
	SubscribeEvent()

	// 处理事件
	HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any)
}

type BaseQuestHandler struct {
}

func (h *BaseQuestHandler) Init(player *playerdomain.Player, quest *playerdomain.Quest) {
	prototype := config.QueryById[configdomain.QuestData](quest.Id)
	if prototype == nil {
		return
	}
	target, _ := util.StringToInt32(prototype.Target)
	quest.Target = target
}

func (h *BaseQuestHandler) GetQuestType() int32 {
	return 0
}


func (h *BaseQuestHandler) OnQuestFinish(player *playerdomain.Player,quest *playerdomain.Quest) {
	 
}

func (h *BaseQuestHandler) CheckProgress(player *playerdomain.Player,quest *playerdomain.Quest) {
	 
}

func (h *BaseQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	 
}

func (h *BaseQuestHandler) SubscribeEvent() {

}

func (h *BaseQuestHandler) Register(handler QuestHandler, topic string) {
	context.EventBus.Subscribe(topic, func(data interface{}) {
		player := data.(events.IPlayerEvent).GetOwner().(*playerdomain.Player)
		if player == nil {
			return
		}
		quests := player.QuestBox.SelectUnFinishedQuestsByType(handler.GetQuestType())
		for _, quest := range quests {
			handler.HandleEvent(player, quest, data)
		}
	})
}

