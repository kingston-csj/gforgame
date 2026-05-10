package config

type RuneData struct {
	Id      int32  `json:"id" excel:"id"`
	Name    string `json:"name" excel:"name"`
	Quality int32  `json:"quality" excel:"quality"`
	Tips    string `json:"tips" excel:"tips"`
	Icon    string `json:"icon" excel:"icon"`
	Overlap int32  `json:"overlap" excel:"overlap"`
}

func (i *RuneData) GetMaxOverlap() int32 {
	return 1
}

func (i *RuneData) GetId() int32 {
	return i.Id
}
