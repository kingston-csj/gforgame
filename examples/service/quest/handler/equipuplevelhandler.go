package handler

import (
	"github.com/forfun/gforgame/examples/constants"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	events "github.com/forfun/gforgame/examples/events"
)

type EquipUpLevelQuestHandler struct {
	BaseQuestHandler
}

func (h *EquipUpLevelQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeEquipUpLevel
}

func (h *EquipUpLevelQuestHandler) SubscribeEvent() {
	h.Register(h, events.EquipLevelUp)
}

func (h *EquipUpLevelQuestHandler) Init(player *playerdomain.Player, quest *playerdomain.Quest) {
	h.BaseQuestHandler.Init(player, quest)
	quest.Progress = player.ExtendBox.EquipUpLevelTimes
}

func (h *EquipUpLevelQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {
	if _, ok := event.(*events.EquipLevelUpEvent); ok {
		quest.Progress++
		h.CheckProgress(player, quest)
	}
}
