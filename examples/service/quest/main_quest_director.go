package quest

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

// 主线任务类别
type MainQuestDirector struct {
	baseQuestDirector
}

func NewMainQuestDirector() *MainQuestDirector {
	return &MainQuestDirector{}
}

/// 实现QuestDirector接口

func (d *MainQuestDirector) OnPlayerLogin(player *playerdomain.Player) {
	firstMainQuestId := config.GetSpecificContainer[container.QuestContainer]("quest").FirstMainQuestId
	if !player.QuestBox.HasReceivedQuest(firstMainQuestId) {
		GetQuestService().AcceptQuest(player, firstMainQuestId)
	}
}

// 玩家完成任务后触发
func (d *MainQuestDirector) AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest) {
	d.AcceptNextQuest(player, quest)
	d.notifyMainQuest(player)
}

func (d *MainQuestDirector) AcceptNextQuest(player *playerdomain.Player, quest *playerdomain.Quest)  {
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	player.QuestBox.AddFinishedQuest(quest.Id)
	if questData.Next != 0 {
		nextQuestData := config.QueryById[configdomain.QuestData](questData.Next)
		GetQuestService().AcceptQuest(player, nextQuestData.Id)
	}
}

// 任务进度变更触发
func (d *MainQuestDirector) OnQuestProgressChanged(player *playerdomain.Player, quest *playerdomain.Quest) {
	questData := config.QueryById[configdomain.QuestData](quest.Id)
	 // 自动领奖
	 if questData.Auto> 0 && quest.IsComplete() {
		// 暂时不处理
	 }
}

func (d *MainQuestDirector) notifyMainQuest(player *playerdomain.Player) {
	quest := player.QuestBox.GetCurrentMainQuest()
	if quest == nil {
		return
	}
	questVo := quest.ToVo()
	refresh := &protos.PushQuestRefreshVo{
		Quest: questVo,
	}
	io.NotifyPlayer(player, refresh)
}

// 获取任务类型
func (d *MainQuestDirector) GetCategoryType() int32 {
	return constants.QuestCategoryMain
}
 







