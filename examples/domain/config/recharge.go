package config

type RechargeData struct {
	Id int32 `json:"id" excel:"id"`

	Money float32 `json:"money" excel:"money"`

	Type int32 `json:"type" excel:"type"`
	// 购买条件，不同的类型有不同的条件
	Condition string `json:"condition" excel:"condition"`

	Rewards string `json:"rewards" excel:"rewards"`
	// 打包售卖关联的孩子id
	Children string `json:"children" excel:"children"`
}