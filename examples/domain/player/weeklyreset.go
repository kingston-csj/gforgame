package player

type WeeklyReset struct {
	// 上一次每周重置的时间戳
	LastWeeklyReset int64
	// 每周任务兑换积分
	WeeklyQuestScore int32
	// 任务周活跃度兑换积分的档位索引（0为未领取，4为全部领取）
	QuestWeeklyRewardIndex int32
	// 充值累计积分
	RechargeSum float32
}

func (d *WeeklyReset) Reset(time int64) {
	d.LastWeeklyReset = time
	d.WeeklyQuestScore = 0
}
