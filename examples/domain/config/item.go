package config

type PropData struct {
	Id              int32  `json:"id" excel:"id"`
	Name            string `json:"name" excel:"name"`
	Quality         int32  `json:"quality" excel:"quality"`
	Tips            string `json:"tips" excel:"tips"`
	Icon            string `json:"icon" excel:"icon"`
	Overlap         int32  `json:"overlap" excel:"overlap"`
	ActivateRewards string `json:"activateRewards" excel:"activateRewards"`
}

func (i *PropData) GetMaxOverlap() int32 {
	return i.Overlap
}

func (i *PropData) GetId() int32 {
	return i.Id
}
