package player

type MonthlyResetBox struct {

	// 重置时间
	ResetTime int64 `json:"reset time"`
	// 签到天数
	SignInDays []int32 `json:"signin days"`

	// 充值累计积分
	RechargeSum float32
}