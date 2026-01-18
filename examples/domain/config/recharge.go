package config

type RechargeData struct {
	Id int32 `json:"id" excel:"id"`

	Money float32 `json:"money" excel:"money"`

	Type int32 `json:"type" excel:"type"`

	Condition string `json:"condition" excel:"condition"`

	Rewards string `json:"rewards" excel:"rewards"`
}