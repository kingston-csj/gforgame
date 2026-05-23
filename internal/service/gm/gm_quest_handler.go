package gm

import (
	commonerrors "github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/common/util/conv"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	questservice "github.com/forfun/gforgame/internal/service/quest"
)

type QuestGmHandler struct {
	quest *questservice.QuestService
}

func NewQuestGmHandler(quest *questservice.QuestService) *QuestGmHandler {
	return &QuestGmHandler{
		quest: quest,
	}
}

func (h *QuestGmHandler) RegisterTo(gm *GmService) {
	gm.Register("quest", "完成任务", "quest 1001", h.handleQuest)
}

func (h *QuestGmHandler) handleQuest(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	questId, _ := conv.StringToInt32(params)
	h.quest.GmFinish(player, questId)
	return nil
}
