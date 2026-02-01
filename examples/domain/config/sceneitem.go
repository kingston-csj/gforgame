package config

type ScenePropData struct {
	PropData
}

func (i *ScenePropData) GetMaxOverlap() int32 {
	return i.Overlap
}

func (i *ScenePropData) GetId() int32 {
	return i.Id
}
