package handler

import (
	"github.com/forfun/gforgame/internal/constants"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	events "github.com/forfun/gforgame/internal/events"
)

type MainGuanKaQuestHandler struct {
	BaseQuestHandler
}

func (h *MainGuanKaQuestHandler) GetQuestType() int32 {
	return constants.QuestTypePassGuanka
}

func (h *MainGuanKaQuestHandler) SubscribeEvent() {
	h.Register(h, events.PassGuanka)
}

func (h *MainGuanKaQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {
	if _, ok := event.(*events.PassGuankaEvent); ok {
		quest.Progress++
		h.CheckProgress(player, quest)
	}
}
