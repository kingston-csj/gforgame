package config

// 任务表
type QuestData struct {
	Id      int32  `json:"id" excel:"id"`
	// 分类 1主线2日常
	Category    int32 `json:"category" excel:"category"`
	// 类型
	Type        int32 `json:"type" excel:"type"`
	// 子类型
	SubType     int32 `json:"subType" excel:"subType"`
	// 任务活跃度
	Score       int32 `json:"score" excel:"score"`
	// 贡献值（公会专用）
	Contribution int32 `json:"contribution" excel:"contribution"`
	// 任务目标，不同类型自行解析
	Target      string `json:"target" excel:"target"`
	// 自动领奖 0手动1自动
	Auto        int32 `json:"auto" excel:"auto"`
	// 前置任务Id (程序动态计算)
	PreviousId int32 `json:"previousId" excel:"previousId"`
	// 下一个任务（主线）
	Next        int32 `json:"next" excel:"next"`
	// 奖励数组字符串
	Rewards      string `json:"rewards" excel:"rewards"`
	// 1继承历史进度（0不继承）
	History int32 `json:"history" excel:"history"`
	// 额外参数 
	Extra string `json:"extra" excel:"extra"`
}

func (q *QuestData) UseHistoryProgress() bool {
	return q.History == 1
}