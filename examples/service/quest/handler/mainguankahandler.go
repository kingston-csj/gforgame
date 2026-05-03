package handler

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
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
