package config

import (
	"io/github/gforgame/domain"
)

type RechargeData struct {
	Id int `json:"id"`

	Money float32 `json:"money"`

	Type string `json:"type"`

	// 购买条件，不同的类型有不同的条件
	Condition string `json:"condition"`

	Rewards []domain.RewardDefLite `json:"rewards"`
}