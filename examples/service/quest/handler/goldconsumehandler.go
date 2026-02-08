package quest

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
)

type GoldConsumeQuestHandler struct {
	BaseQuestHandler
}

func (h *GoldConsumeQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeGoldConsume
}

func (h *GoldConsumeQuestHandler) SubscribeEvent() {
	h.Register(h, events.ItemConsume)
}

func (h *GoldConsumeQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	if evt, ok := event.(*events.ItemConsumeEvent); ok {
		if evt.ItemId == constants.ITEM_GOLD_ID {
			quest.Progress += evt.Count
			h.CheckProgress(player, quest)
		}
	}
}