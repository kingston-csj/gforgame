package config

// 每日签到奖励
type SigninData struct {
	Id      int32  `json:"id" excel:"id"`
	Rewards string `json:"rewards" excel:"rewards"`
}
