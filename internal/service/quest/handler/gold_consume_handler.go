package handler

import (
	"github.com/forfun/gforgame/internal/constants"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	events "github.com/forfun/gforgame/internal/events"
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

func (h *GoldConsumeQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {
	if evt, ok := event.(*events.ItemConsumeEvent); ok {
		if evt.ItemId == constants.ITEM_GOLD_ID {
			quest.Progress += evt.Count
			h.CheckProgress(player, quest)
		}
	}
}
