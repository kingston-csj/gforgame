package handler

import (
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/context"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	events "github.com/forfun/gforgame/internal/events"
	qcore "github.com/forfun/gforgame/internal/service/quest/core"
)

type BaseQuestHandler struct {
	resolver qcore.Resolver
}

func (h *BaseQuestHandler) Init(player *playerdomain.Player, quest *playerdomain.Quest) {
	prototype := config.QueryById[configdomain.QuestData](quest.Id)
	if prototype == nil {
		return
	}
	quest.Target = prototype.Target
}

func (h *BaseQuestHandler) GetQuestType() int32 {
	return 0
}

func (h *BaseQuestHandler) OnQuestFinish(player *playerdomain.Player, quest *playerdomain.Quest) {

}

func (h *BaseQuestHandler) CheckProgress(player *playerdomain.Player, quest *playerdomain.Quest) {
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	if questData == nil || h.resolver == nil {
		return
	}
	questDirector := h.resolver.GetQuestDirector(questData.Category)
	if questDirector == nil {
		return
	}
	questDirector.OnQuestProgressChanged(player, quest)
}

func (h *BaseQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {

}

func (h *BaseQuestHandler) SubscribeEvent() {

}

func (h *BaseQuestHandler) SetResolver(resolver qcore.Resolver) {
	h.resolver = resolver
}

func (h *BaseQuestHandler) Register(handler qcore.QuestHandler, topic string) {
	context.EventBus.Subscribe(topic, func(data interface{}) {
		var player *playerdomain.Player
		switch v := data.(type) {
		case *playerdomain.Player:
			player = v
		case events.IPlayerEvent:
			typedPlayer, ok := v.GetOwner().(*playerdomain.Player)
			if ok {
				player = typedPlayer
			}
		}
		if player == nil {
			return
		}
		quests := player.QuestBox.SelectUnFinishedQuestsByType(handler.GetQuestType())
		for _, quest := range quests {
			handler.HandleEvent(player, quest, data)
		}
	})
}
