package protos

type ActivityRewardVo struct {
	Id    int32  `json:"id"`    // 关联activityreward表id
	Value string `json:"value"` //奖励信息（可拓展） 0未领奖，1已完成未领奖，2已领奖
}

type ActivityVo struct {
	ActivityId int32               `json:"activityId"` // 关联activity表id
	RewardVos  []*ActivityRewardVo `json:"rewardVos"`  // 奖励信息
}

type PushActivityLoadAll struct {
	_           struct{}      `cmd_ref:"CmdPushActivityLoadAll"`
	ActivityVos []*ActivityVo `json:"activityVos"` // 活动信息
}
