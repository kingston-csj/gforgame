package handler

import (
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/constants"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	events "github.com/forfun/gforgame/internal/events"
)

type LoginQuestHandler struct {
	BaseQuestHandler
}

func (h *LoginQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeLogin
}

func (h *LoginQuestHandler) SubscribeEvent() {
	h.Register(h, events.PlayerLogin)
}

func (h *LoginQuestHandler) HandleEvent(player *playerdomain.Player, quest *playerdomain.Quest, event any) {
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	if questData.UseHistoryProgress() {
		quest.SetProgress(player.ExtendBox.AccumulatedLoginDays)
	} else {
		quest.AddProgress(1)
	}
	h.CheckProgress(player, quest)
}
