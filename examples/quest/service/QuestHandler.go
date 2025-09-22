package quest

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/network"
)
type QuestHandler interface {

	Init();
	// 对应的任务类型
	GetQuestType()int32;
	// 执行完成时触发
	OnQuestFinish(player *playerdomain.Player, quest *playerdomain.Quest);
	// 检查任务进度
	CheckProgress(player *playerdomain.Player, quest *playerdomain.Quest);
}

type BaseQuestHandler struct {
}

func (h *BaseQuestHandler) GetQuestType() int32 {
	return 0
}

func (h *BaseQuestHandler) OnQuestFinish(player *playerdomain.Player, quest *playerdomain.Quest) {
	quest.Status = constants.QuestStatusFinished
}

func (h *BaseQuestHandler) CheckProgress(player *playerdomain.Player, quest *playerdomain.Quest) {
	if quest.IsComplete() {
		h.OnQuestFinish(player, quest)
	}
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	if questData == nil {
		panic("quest data not found")
	}
	director := GetQuestDirector(questData.Category)
	director.OnQuestProgressChanged(player, quest)
}

type RecruitQuestHandler struct {
	BaseQuestHandler
}

func (h *RecruitQuestHandler) GetQuestType() int32 {
	return constants.QuestTypeRecruit
}

func (h *RecruitQuestHandler) Init() {
	context.EventBus.Subscribe(events.Recruit, func(data interface{}) {
		event := data.(*events.RecruitEvent)
		p:=network.GetPlayerByPlayerId(event.Player.GetId())
		if p == nil {
			return
		}
		player := p.(*playerdomain.Player)
		quests:=player.QuestBox.SelectUnFinishedQuestsByType(h.GetQuestType())
		for _, quest := range quests {
			quest.AddProgress(1)
			h.CheckProgress(player, quest)
		}
	})
}