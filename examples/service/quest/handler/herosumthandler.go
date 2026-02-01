package quest

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
)

type HeroSumQuestHandler struct {
	BaseQuestHandler
}

func (h *HeroSumQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeHeroSum
}

func (h *HeroSumQuestHandler) SubscribeEvent() {
	h.Register(h, events.HeroGain)
}

func (h *HeroSumQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	if _, ok := event.(*events.HeroGainEvent); ok {
		quest.SetProgress(int32(len(player.HeroBox.GetAllHeros())))
		h.CheckProgress(player, quest)
	}
}