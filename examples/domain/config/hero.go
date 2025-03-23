package config

type HeroData struct {
	Id      int32  `json:"id" excel:"id"`
	Name    string `json:"name" excel:"name"`
	Quality int32  `json:"quality" excel:"quality"`
	Tips    string `json:"tips" excel:"tips"`
	Icon    string `json:"icon" excel:"icon"`
	// 抽奖概率
	Prob int32 `json:"prob" excel:"prob"`
	// 对应的碎片数量
	Shard int32 `json:"shard" excel:"shard"`
	// 对应的碎片道具id
	Item int32 `json:"item" excel:"item"`
}
