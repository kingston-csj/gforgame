package player

type DailyReset struct {
	// 上一次每日重置的时间戳
	LastDailyReset int64
	// 每日任务兑换积分
	DailyQuestScore int32
	// 任务日活跃度兑换积分的档位索引（0为未领取，4为全部领取）
	QuestDailyRewardIndex int32
}

func (d *DailyReset) Reset(time int64) {
	d.LastDailyReset = time
	d.DailyQuestScore = 0
}
