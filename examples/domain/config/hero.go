package config

import (
	"io/github/gforgame/examples/fight/attribute"
)

type HeroData struct {
	Id      int32  `json:"id" excel:"id"`
	Name    string `json:"name" excel:"name"`
	Quality int32  `json:"quality" excel:"quality"`
	Tips    string `json:"tips" excel:"tips"`
	Icon    string `json:"icon" excel:"icon"`
	// 技能
	Skills []int32 `json:"skills" excel:"skills"`
	// 抽奖概率
	Prob int32 `json:"prob" excel:"prob"`
	// 对应的碎片数量
	Shard int32 `json:"shard" excel:"shard"`
	// 对应的碎片道具id
	Item int32 `json:"item" excel:"item"`
	// 生命值
	Hp int32 `json:"hp"	 excel:"hp"`
	// 攻击力
	Attack int32 `json:"attack" excel:"attack"`
	// 防御力
	Defense int32 `json:"defense" excel:"defense"`
	// 速度
	Speed int32 `json:"speed" excel:"speed"`
	// 属性（手动收集）
	Attrs []attribute.Attribute `json:"-" excel:"-"`
}

func (h *HeroData) GetHeroAttrs() []attribute.Attribute {
	if h.Attrs == nil {
		h.Attrs = make([]attribute.Attribute, 4)
		h.Attrs[0] = attribute.Attribute{
			AttrType: attribute.Hp,
			Value:    int32(h.Hp),
		}
		h.Attrs[1] = attribute.Attribute{
			AttrType: attribute.Attack,
			Value:    int32(h.Attack),
		}
		h.Attrs[2] = attribute.Attribute{
			AttrType: attribute.Defense,
			Value:    int32(h.Defense),
		}
		h.Attrs[3] = attribute.Attribute{
			AttrType: attribute.Speed,
			Value:    int32(h.Speed),
		}
	}
	return h.Attrs
}
