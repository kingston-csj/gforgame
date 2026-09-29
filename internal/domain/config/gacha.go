package config

import domainreward "github.com/forfun/gforgame/internal/domain/reward"

type GachaData struct {
	Id int32 `json:"id" excel:"id"`
	// 类型 1为普通招募，2为高级招募
	Type int32 `json:"type" excel:"type"`
	// 权重
	Weight int32 `json:"weight" excel:"weight"`
	// 奖励
	Rewards domainreward.Reward `json:"rewards" excel:"rewards"`
}
