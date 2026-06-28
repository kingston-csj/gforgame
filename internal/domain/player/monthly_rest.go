package player

type MonthlyResetBox struct {

	// 重置时间
	ResetTime int64 `json:"reset time"`
	// 签到天数
	SignInDays []int32 `json:"signin days"`
	// 充值累计积分
	RechargeSum float32
	// 补签：目标天数对应补签天数
	SignInMakeUp map[int32]int32 `json:"signInMakeUp"`
}

func (m *MonthlyResetBox) AfterLoad() {
	if m.SignInDays == nil {
		m.SignInDays = make([]int32, 0)
	}
	if m.SignInMakeUp == nil {
		m.SignInMakeUp = make(map[int32]int32)
	}
}