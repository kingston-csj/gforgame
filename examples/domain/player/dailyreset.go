package player

type DailyReset struct {
	// 上一次每日重置的时间戳
	LastDailyReset int64
	// 每日任务兑换积分
	DailyQuestScore int32
}

func (d *DailyReset) Reset(time int64) {
	d.LastDailyReset = time
	d.DailyQuestScore = 0
}
