package quest

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
)

type TimeQuestHandler struct {
	BaseQuestHandler
}

func (h *TimeQuestHandler) GetQuestType() int32 {
	return constants.QuestTime
}

func (h *TimeQuestHandler) SubscribeEvent() {
}


func (h *TimeQuestHandler) HandleEvent(player *playerdomain.Player,quest *playerdomain.Quest, event any) {
	 // 进度由定时器触发
}