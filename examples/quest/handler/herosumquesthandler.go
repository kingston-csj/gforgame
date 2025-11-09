package quest

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
)

type HeroSumQuestHandler struct {
	BaseQuestHandler
}

func (h *HeroSumQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	if _, ok := event.(*events.HeroGainEvent); ok {
		quest.SetProgress(int32(len(player.HeroBox.GetAllHeros())))
		h.CheckProgress(player, quest)
	}
}

func (h *HeroSumQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeHeroSum
}

func (h *HeroSumQuestHandler) SubscribeEvent() {
	context.EventBus.Subscribe(events.HeroGain, func(data interface{}) {
		event := data.(*events.HeroGainEvent)
		p := event.Player.(*playerdomain.Player)
		quests := p.QuestBox.SelectUnFinishedQuestsByType(h.GetQuestType())
		for _, q := range quests {
			h.HandleEvent(p, q, data)
		}
	})
}
