package config

// ActivityRewardData 活动奖励配置表
type ActivityRewardData struct {
	Id    int32  `json:"id" excel:"id"`
	// ActivityId 活动ID
	ActivityId int32 `json:"activityId" excel:"activityId"`
	// Rewards 奖励
	Rewards  string `json:"rewards" excel:"rewards"`
	// Condition 条件
	Condition  string `json:"condition" excel:"condition"`
}