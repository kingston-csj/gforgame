package director

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/protos"
	qcore "io/github/gforgame/examples/service/quest/core"
)

type baseQuestDirector struct {
	self qcore.QuestDirector
	resolver qcore.Resolver
}

func NewBaseQuestDirector(self qcore.QuestDirector) *baseQuestDirector {
	return &baseQuestDirector{
		self: self,
	}
}



func (d *baseQuestDirector) AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest) {

}

// 任务进度变更触发
func (d *baseQuestDirector) OnQuestProgressFinished(player *playerdomain.Player, quest *playerdomain.Quest) {
	quest.Status = constants.QuestStatusFinished
}

func (d *baseQuestDirector) TakeProgressRewards(player *playerdomain.Player) []*protos.RewardVo {
	return make([]*protos.RewardVo, 0)
}

func (d *baseQuestDirector) OnQuestProgressChanged(player *playerdomain.Player, quest *playerdomain.Quest) {
	if quest.IsComplete() {
		d.self.OnQuestProgressFinished(player, quest)
	}
	questVo := quest.ToVo()
	refresh := &protos.PushQuestRefreshVo{
		Quest: questVo,
	}
	io.NotifyPlayer(player, refresh)
}

func (d *baseQuestDirector) SetResolver(resolver qcore.Resolver) {
	d.resolver = resolver
}
