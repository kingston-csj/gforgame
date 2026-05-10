package config

type VipData struct {
	Id int32 `json:"id"`

	Money int32 `json:"money"`

	Type string `json:"type"`

	Rewards string `json:"rewards"`
	// 竞技场额外次数
	ArenaTimes int32 `json:"arenaTimes"`
}