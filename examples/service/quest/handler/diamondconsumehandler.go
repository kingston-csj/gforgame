package quest

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
)

type DiamondConsumeQuestHandler struct {
	BaseQuestHandler
}

func (h *DiamondConsumeQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeDiamondConsume
}

func (h *DiamondConsumeQuestHandler) SubscribeEvent() {
	h.Register(h, events.ItemConsume)
}

func (h *DiamondConsumeQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	if evt, ok := event.(*events.ItemConsumeEvent); ok {
		if evt.ItemId == constants.ITEM_DIAMOND_ID {
			quest.Progress += evt.Count
			h.CheckProgress(player, quest)
		}
	}
}