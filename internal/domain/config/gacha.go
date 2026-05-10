package config

type GachaData struct {
	Id int32 `json:"id" excel:"id"`
	// 类型 1为普通招募，2为高级招募
	Type int32 `json:"type" excel:"type"`
	// 权重
	Weight int32 `json:"weight" excel:"weight"`
	// 奖励
	Rewards string `json:"rewards" excel:"rewards"`
}