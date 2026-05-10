package config

// ActivityData 活动配置表
type ActivityData struct {
	Id    int32  `json:"id" excel:"id"`
	// Start 活动开始时间
	Start string `json:"start" excel:"start"`
	// End 活动结束时间
	End   string `json:"end" excel:"end"`
	// Cron 其他自定义调度表达式
	Cron  string `json:"cron" excel:"cron"`
}