package config

type MailData struct {
	Id int32 `json:"id" excel:"id"`
	// 有效时间（小时）
	ValidTime int32  `json:"validTime" excel:"validTime"`
	Rewards   string `json:"rewards" excel:"rewards"`
}