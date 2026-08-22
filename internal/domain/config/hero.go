package config

import (
	"github.com/forfun/gforgame/internal/domain/attr"
)

type HeroData struct {
	// 属性
	attr.AttributeItems `json:"-" excel:"-"`
	Id                  int32  `json:"id" excel:"id"`
	Name                string `json:"name" excel:"name"`
	Quality             int32  `json:"quality" excel:"quality"`
	Tips                string `json:"tips" excel:"tips"`
	Icon                string `json:"icon" excel:"icon"`
	// 技能
	Skills []int32 `json:"skills" excel:"skills"`
	// 对应的碎片数量
	ShardAmount int32 `json:"shard" excel:"shardAmount"`
	// 对应的碎片道具id
	ShardItem int32 `json:"shardItem" excel:"shardItem"`
	// 生命值
	Hp int32 `json:"hp"	 excel:"hp"`
	// 攻击力
	Attack int32 `json:"attack" excel:"attack"`
	// 防御力
	Defense int32 `json:"defense" excel:"defense"`
	// 速度
	Speed int32 `json:"speed" excel:"speed"`
}
