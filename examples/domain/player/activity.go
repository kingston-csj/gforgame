package player

type ActivityInfo struct {
	Rewards map[int32]string `json:"rewards"`
}

type ActivityBox struct {
	Data map[int32]*ActivityInfo `json:"data"`
}

func (b *ActivityBox) AfterLoad() {
	if b.Data == nil {
		b.Data = make(map[int32]*ActivityInfo)
	}
}
