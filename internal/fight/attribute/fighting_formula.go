package attribute

import "github.com/forfun/gforgame/internal/domain/attr"

func CalculateFightingPower(attrBox *attr.AttrBox) int32 {
	power := int32(0)
	attrs := attrBox.GetAttrs()
	for _, attr := range attrs {
		power += int32(attr.Value)
	}
	return power
}
