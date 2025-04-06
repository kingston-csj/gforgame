package attribute

type AttrBox struct {
	Attrs map[AttrType]float32 `json:"-"`
}

func NewAttrBox() *AttrBox {
	return &AttrBox{
		Attrs: make(map[AttrType]float32),
	}
}

func (a *AttrBox) GetAttr(attrType AttrType) Attribute {
	return Attribute{
		AttrType: attrType,
		Value:    a.Attrs[attrType],
	}
}

func (a *AttrBox) AddAttr(attrType AttrType, value float32) {
	a.Attrs[attrType] += value
}

func (a *AttrBox) AddAttrs(attrs []Attribute) {
	for _, attr := range attrs {
		a.AddAttr(attr.AttrType, attr.Value)
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
