package quest

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
)

type HeroLevelUpQuestHandler struct {
	BaseQuestHandler
}

func (h *HeroLevelUpQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeHeroUpLevel
}

func (h *HeroLevelUpQuestHandler) SubscribeEvent() {
	h.Register(h, events.HeroLevelUp)
}

func (h *HeroLevelUpQuestHandler) Init(player *playerdomain.Player,quest *playerdomain.Quest) {
	quest.Progress = player.HeroBox.UpLevelTimes
}

func (h *HeroLevelUpQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	if _, ok := event.(*events.HeroLevelUpEvent); ok {
		quest.Progress = player.HeroBox.UpLevelTimes
		h.CheckProgress(player, quest)
	}
}