package config

import "io/github/gforgame/examples/attribute"

type HeroStageData struct {
	Id       int32 `json:"id" excel:"id"`
	MaxLevel int32 `json:"max_level" excel:"max_level"`
	Cost     int32 `json:"cost" excel:"cost"`
	Hp       int32 `json:"hp" excel:"hp"`
	Attack   int32 `json:"attack" excel:"attack"`
	Defense  int32 `json:"defense" excel:"defense"`
	Speed    int32 `json:"speed" excel:"speed"`
	// 属性（手动收集）
	Attrs []attribute.Attribute `json:"-" excel:"-"`
}

func (h *HeroLevelData) GetHeroStageAttrs() []attribute.Attribute {
	if h.Attrs == nil {
		h.Attrs = make([]attribute.Attribute, 4)
		h.Attrs[0] = attribute.Attribute{
			AttrType: attribute.Hp,
			Value:    float32(h.Hp),
		}
		h.Attrs[1] = attribute.Attribute{
			AttrType: attribute.Attack,
			Value:    float32(h.Attack),
		}
		h.Attrs[2] = attribute.Attribute{
			AttrType: attribute.Defense,
			Value:    float32(h.Defense),
		}
		h.Attrs[3] = attribute.Attribute{
			AttrType: attribute.Speed,
			Value:    float32(h.Speed),
		}
	}
	return h.Attrs
}
