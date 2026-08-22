package attr

type Attribute struct {
	AttrType AttrType
	Value    int32
}

type AttributeItems struct {
	Hp               int32 `json:"hp" excel:"hp"`
	Attack           int32 `json:"attack" excel:"attack"`
	Defense          int32 `json:"defense" excel:"defense"`
	Speed            int32 `json:"speed" excel:"speed"`
	Move_speed       int32 `json:"move_speed" excel:"move_speed"`
	Attack_range     int32 `json:"attack_range" excel:"attack_range"`
	Attack_speed_add int32 `json:"attack_speed_add" excel:"attack_speed_add"`
	Cd               int32 `json:"cd" excel:"cd"`
	Hp_add           int32 `json:"hp_add" excel:"hp_add"`
	Attack_add       int32 `json:"attack_add" excel:"attack_add"`
	Defense_add      int32 `json:"defense_add" excel:"defense_add"`
	Crit_probability int32 `json:"crit_probability" excel:"crit_probability"`
	Crit_hit         int32 `json:"crit_hit" excel:"crit_hit"`
	Penetration      int32 `json:"penetration" excel:"penetration"`
	// 属性（手动收集）
	Attrs []Attribute `json:"-" excel:"-"`
}

func (h *AttributeItems) GetAttrs() []Attribute {
	if h.Attrs != nil {
		return h.Attrs
	}
	var attrs []Attribute
	h.tryAdd(&attrs, Hp, h.Hp)
	h.tryAdd(&attrs, Attack, h.Attack)
	h.tryAdd(&attrs, Defense, h.Defense)
	h.tryAdd(&attrs, Speed, h.Speed)
	h.tryAdd(&attrs, Move_speed, h.Move_speed)
	h.tryAdd(&attrs, Attack_range, h.Attack_range)
	h.tryAdd(&attrs, Attack_speed_add, h.Attack_speed_add)
	h.tryAdd(&attrs, Cd, h.Cd)
	h.tryAdd(&attrs, Hp_add, h.Hp_add)
	h.tryAdd(&attrs, Attack_add, h.Attack_add)
	h.tryAdd(&attrs, Defense_add, h.Defense_add)
	h.tryAdd(&attrs, Crit_probability, h.Crit_probability)
	h.tryAdd(&attrs, Crit_hit, h.Crit_hit)
	h.tryAdd(&attrs, Penetration, h.Penetration)
	return attrs
}

// tryAdd 数值非0时才追加属性
func (h *AttributeItems) tryAdd(list *[]Attribute, typ AttrType, val int32) {
	if val == 0 {
		return
	}
	*list = append(*list, Attribute{
		AttrType: typ,
		Value:    val,
	})
}
