package quest

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
)

type ClientEventQuestHandler struct {
	BaseQuestHandler
}

func (h *ClientEventQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeClientEvent
}

func (h *ClientEventQuestHandler) SubscribeEvent() {
	h.Register(h, events.ClientDiyEvent)
}


func (h *ClientEventQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	if evt, ok := event.(*events.ClientCustomEvent); ok {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		if questData.SubType <= 0 {
			return
		}
		eventId := evt.EventId
		if questData.UseHistoryProgress() {
			quest.SetProgress(player.ExtendBox.ClientEvents[eventId])
		} else {
			quest.AddProgress(1)
		}
		h.CheckProgress(player, quest)
	}
}