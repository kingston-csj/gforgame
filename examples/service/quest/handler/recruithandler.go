package handler

import (
	"github.com/forfun/gforgame/examples/config"
	"github.com/forfun/gforgame/examples/constants"
	configdomain "github.com/forfun/gforgame/examples/domain/config"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	events "github.com/forfun/gforgame/examples/events"
)

type RecruitQuestHandler struct {
	BaseQuestHandler
}

func (h *RecruitQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeRecruit
}

func (h *RecruitQuestHandler) SubscribeEvent() {
	h.Register(h, events.Recruit)
}

func (h *RecruitQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {
	if evt, ok := event.(*events.RecruitEvent); ok {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		if questData.SubType > 0 {
			if questData.SubType != evt.Type {
				return
			}
		}
		quest.AddProgress(1)
		h.CheckProgress(player, quest)
	}
}
