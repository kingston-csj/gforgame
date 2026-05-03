package handler

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
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
