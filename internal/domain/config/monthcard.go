package config

type MonthlyCardData struct {
	Id      int32  `json:"id" excel:"id"`
	Rewards string `json:"rewards" excel:"rewards"`
	// 竞技场额外挑战次数
	ArenaTimes int32 `json:"arenaTimes" excel:"arenaTimes"`
	// 副本额外次数
	FubenTimes int32 `json:"fubenTimes" excel:"fubenTimes"`
	// 挂机额外时间（分钟）
	IdleAddMinutes int32 `json:"idleAddMinutes" excel:"idleAddMinutes"`
	// 有效天数
	ValidDays int32 `json:"validDays" excel:"validDays"`
}