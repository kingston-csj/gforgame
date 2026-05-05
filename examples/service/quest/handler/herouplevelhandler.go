package handler

import (
	"github.com/forfun/gforgame/examples/constants"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	events "github.com/forfun/gforgame/examples/events"
)

type HeroLevelUpQuestHandler struct {
	BaseQuestHandler
}

func (h *HeroLevelUpQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeHeroUpLevel
}

func (h *HeroLevelUpQuestHandler) SubscribeEvent() {
	h.Register(h, events.HeroLevelUp)
}

func (h *HeroLevelUpQuestHandler) Init(player *playerdomain.Player, quest *playerdomain.Quest) {
	h.BaseQuestHandler.Init(player, quest)
	quest.Progress = player.HeroBox.UpLevelTimes
}

func (h *HeroLevelUpQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {
	if _, ok := event.(*events.HeroLevelUpEvent); ok {
		quest.Progress = player.HeroBox.UpLevelTimes
		h.CheckProgress(player, quest)
	}
}
