package handler

import (
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/constants"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	events "github.com/forfun/gforgame/internal/events"
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
