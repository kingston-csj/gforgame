package attribute

func CalculateFightingPower(attrBox *AttrBox) int32 {
	power := int32(0)
	attrs := attrBox.GetAttrs()
	for _, attr := range attrs {
		power += int32(attr.Value)
	}
	return power
}
