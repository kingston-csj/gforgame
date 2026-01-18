package player

type RechargeBox struct {
	// 七日活动：是否已激活付费奖励
	ActivatedQiRiPay int32 `json:"activatedQiRiPay"`
	// 充值累计次数
	RechargeTimes map[int32]int32 `json:"rechargeTimes"`
	// 通行证活动：是否已激活付费奖励
	ActivatedPassPay int32 `json:"activatedPassPay"`
	// 首充时间
	FirstRechargeTime int64 `json:"firstRechargeTime"`
}