package quest

import (
	playerdomain "io/github/gforgame/examples/domain/player"
)

// 任务类型切面控制器
type QuestDirector interface {
    //  玩家登录触发，下发任务信息
    OnPlayerLogin(player *playerdomain.Player)

    // 玩家完成任务后触发
    AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest)

    // 任务进度变更触发
    OnQuestProgressChanged(player *playerdomain.Player, quest *playerdomain.Quest)

	// 获取任务类型
	GetCategoryType(quest *playerdomain.Quest) int
}