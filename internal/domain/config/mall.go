package config

// 每日签到奖励商城数据
type MallData struct {
	Id      int32  `json:"id" excel:"id"`
	Type    int32  `json:"type" excel:"type"`
	Rewards string `json:"rewards" excel:"rewards"`
	Consumes string `json:"consume" excel:"consumes"`
	// 每日限购
	DailyBuy int32 `json:"dailyBuy" excel:"dailyBuy"`
	// 终身限购
	LifeTimeBuy int32 `json:"lifeTimeBuy" excel:"lifeTimeBuy"`
}
