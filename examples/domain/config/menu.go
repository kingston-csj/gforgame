package config

type MenuData struct {
	PropData
}

func (i *MenuData) GetMaxOverlap() int32 {
	return i.Overlap
}

func (i *MenuData) GetId() int32 {
	return i.Id
}
