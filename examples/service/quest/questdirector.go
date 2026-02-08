package quest

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

// 任务类型切面控制器
type QuestDirector interface {
    //  玩家登录触发，下发任务信息
    OnPlayerLogin(player *playerdomain.Player)

    // 玩家完成任务后触发
    AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest)

    // 任务进度变更触发
    OnQuestProgressChanged(player *playerdomain.Player, quest *playerdomain.Quest)

    // 任务完成执行切面
    OnQuestProgressFinished(player *playerdomain.Player, quest *playerdomain.Quest)

    // 领取任务进度奖励
    TakeProgressRewards(player *playerdomain.Player) []*protos.RewardVo

	// 获取任务类型
	GetCategoryType() int32
}

type baseQuestDirector struct {
}

func (d *baseQuestDirector) AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest) {

}

// 任务进度变更触发
func (d *baseQuestDirector) OnQuestProgressFinished(player *playerdomain.Player, quest *playerdomain.Quest) {

}

func (d *baseQuestDirector) TakeProgressRewards(player *playerdomain.Player) []*protos.RewardVo{
    return make([]*protos.RewardVo, 0)
}

func (d *baseQuestDirector) OnQuestProgressChanged(player *playerdomain.Player, quest *playerdomain.Quest) {
    if quest.IsComplete() {
        d.OnQuestProgressFinished(player, quest)
    }
    questVo := quest.ToVo()
    refresh := &protos.PushQuestRefreshVo{
        Quest:  questVo,
    }
   io.NotifyPlayer(player, refresh)
}