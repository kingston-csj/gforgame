package config

// 任务表
type QuestData struct {
	Id      int32  `json:"id" excel:"id"`
	// 分类 1主线2日常
	Category    int32 `json:"category" excel:"category"`
	// 子类型
	Type        int32 `json:"type" excel:"type"`
	// 任务活跃度
	Score       int32 `json:"score" excel:"score"`
	// 贡献值（公会专用）
	Contribution int32 `json:"contribution" excel:"contribution"`
	// 任务目标，不同类型自行解析
	Target      string `json:"target" excel:"target"`
	// 自动领奖 0手动1自动
	Auto        int32 `json:"auto" excel:"auto"`
	// 下一个任务（主线）
	Next        int32 `json:"next" excel:"next"`
	// 奖励数组字符串
	Rewards      string `json:"rewards" excel:"rewards"`
}