package attribute

type AttrBox struct {
	Attrs map[AttrType]int32 `json:"-"`
}

func NewAttrBox() *AttrBox {
	return &AttrBox{
		Attrs: make(map[AttrType]int32),
	}
}

func (a *AttrBox) GetAttr(attrType AttrType) Attribute {
	return Attribute{
		AttrType: attrType,
		Value:    a.Attrs[attrType],
	}
}

func (a *AttrBox) AddAttr(attrType AttrType, value int32) {
	a.Attrs[attrType] += value
}

func (a *AttrBox) AddAttrs(attrs []Attribute) {
	for _, attr := range attrs {
		a.AddAttr(attr.AttrType, int32(attr.Value))
	}
}

func (a *AttrBox) GetAttrs() []Attribute {
	attrs := make([]Attribute, 0, len(a.Attrs))
	for attrType, value := range a.Attrs {
		attrs = append(attrs, Attribute{
			AttrType: attrType,
			Value:    value,
		})
	}
	return attrs
}

func (a *AttrBox) GetAttrValue(attrType AttrType) int32 {
	v, ok := a.Attrs[attrType]
	if !ok {
		return 0
	}
	return int32(v)
}
