package handler

import (
	"github.com/forfun/gforgame/examples/constants"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	events "github.com/forfun/gforgame/examples/events"
)

type FubenLevelQuestHandler struct {
	BaseQuestHandler
}

func (h *FubenLevelQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeFuben
}

func (h *FubenLevelQuestHandler) SubscribeEvent() {
	h.Register(h, events.PassFuben)
}

func (h *FubenLevelQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {
	if _, ok := event.(*events.PassFubenEvent); ok {
		quest.Progress++
		h.CheckProgress(player, quest)
	}
}
