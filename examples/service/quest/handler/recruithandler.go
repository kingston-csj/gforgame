package quest

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
	"io/github/gforgame/network"
)

type RecruitQuestHandler struct {
	BaseQuestHandler
}

func (h *RecruitQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeRecruit
}

func (h *RecruitQuestHandler) SubscribeEvent() {
	context.EventBus.Subscribe(events.Recruit, func(data interface{}) {
		event := data.(*events.RecruitEvent)
		p := network.GetPlayerByPlayerId(event.Player.GetId())
		if p == nil {
			return
		}
		player := p.(*playerdomain.Player)
		quests := player.QuestBox.SelectUnFinishedQuestsByType(h.GetQuestType())
		for _, quest := range quests {
			h.HandleEvent(player, quest, data)
		}
	})
}