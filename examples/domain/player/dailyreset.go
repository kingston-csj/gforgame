package player

type DailyReset struct {
	// 上一次每日重置的时间戳
	LastDailyReset int64
	// 每日任务兑换积分
	DailyQuestScore int32
	// 任务日活跃度兑换积分的档位索引（0为未领取，4为全部领取）
	QuestDailyRewardIndex int32
	// 充值累计积分
	RechargeSum float32
	// 银月卡奖励是否已领取
	SilverMonthCardReward bool
	// 金月卡奖励是否已领取
	GoldMonthCardReward bool
	// 普通招募免费次数是否已使用
	NormalRecruitFreeUsed bool
	// 普通招募免费次数
	NormalRecruitTimes int32
	// 高级卡免费招募次数是否已使用
	HighRecruitFreeUsed bool
	// 高级卡免费招募次数
	HighRecruitTimes int32
}

func (d *DailyReset) Reset(time int64) {
	d.LastDailyReset = time
	d.DailyQuestScore = 0
}
