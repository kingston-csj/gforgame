package protos

type QuestVo struct {
	Id int32 `json:"id"`
	// 进度
	Progress int32 `json:"progress"`
	// 目标
	Target int32 `json:"target"`
	// 状态 0未完成，1已完成未领奖,2已领奖
	Status int8 `json:"status"`
}

// 主线任务——自动领奖
type PushQuestAutoTakeReward struct {
	_         struct{}    `cmd_ref:"CMD_PUSH_QUEST_AUTO_REWARD" type:"push"`
	RewardVos []*RewardVo `json:"rewardVos"`
}

type PushQuestDailyInfo struct {
	_ struct{} `cmd_ref:"CMD_PUSH_DAILY_QUEST" type:"push"`
	// 已领取的档位索引（0为未领取)
	DailyRewardIndex int32 `json:"dailyRewardIndex"`
	// 今日活跃度
	DailyScore int32 `json:"dailyScore"`
	// 所有任务
	Quests []*QuestVo `json:"quests"`
}

type PushQuestRefreshVo struct {
	_     struct{} `cmd_ref:"CMD_PUSH_UPDATE_QUEST" type:"push"`
	Quest *QuestVo `json:"quest"`
}

type PushQuestReplace struct {
	_ struct{} `cmd_ref:"CMD_RES_REPLACE_QUEST" type:"push"`
	// 旧任务id
	OldQuestId int32 `json:"oldQuestId"`

	Quest *QuestVo `json:"quest"`
}

// 每周任务主界面信息
type PushQuestWeeklyInfo struct {
	_ struct{} `cmd_ref:"CMD_PUSH_WEEKLY_QUEST" type:"push"`
	// 已领取的档位索引（0为未领取)
	WeeklyRewardIndex int32 `json:"weeklyRewardIndex"`
	// 本周活跃度
	WeeklyScore int32 `json:"weeklyScore"`
	// 所有任务
	Quests []*QuestVo `json:"quests"`
}

// 任务——一键领取所有奖励
type ReqQuestTakeAllRewards struct {
	_ struct{} `cmd_ref:"CMD_REQ_QUEST_ALL_REWARD" type:"req"`

	// 任务类型 1主线，2日常
	Category int32 `json:"category"`
}

// 任务——领取档位奖
type ReqQuestTakeProgressReward struct {
	_ struct{} `cmd_ref:"CMD_REQ_QUEST_PROGRESS_REWARD" type:"req"`

	// 任务类型 2每日，5每周，6公会
	Type int32 `json:"type"`
}

// 任务——领取达标奖
type ReqQuestTakeReward struct {
	_ struct{} `cmd_ref:"CMD_REQ_QUEST_REWARD" type:"req"`

	// 任务id
	Id int32 `json:"id"`
}

// 任务——一键领取所有奖励
type ResQuestTakeAllRewards struct {
	_ struct{} `cmd_ref:"CMD_RES_QUEST_ALL_REWARD" type:"res"`

	// 奖励vo
	RewardVos []*RewardVo `json:"rewardVos"`

	// 今日活跃度
	DailyScore int32 `json:"dailyScore"`

	// 本周活跃度
	WeeklyScore int32 `json:"weeklyScore"`

	// 已领取的任务id列表
	QuestIds []int32 `json:"questIds"`
	// 总分数
	Score int32 `json:"score"`
}

type ResQuestTakeProgressReward struct {
	_ struct{} `cmd_ref:"CMD_RES_QUEST_PROGRESS_REWARD" type:"res"`
	// 任务类型 2每日，5每周，6公会
	Type int32 `json:"type"`
	// 已领取的档位索引（0为未领取)
	RewardIndex int32 `json:"rewardIndex"`
	// 奖励vo
	RewardVos []*RewardVo `json:"rewardVos"`
}

// 任务——领取达标奖
type ResQuestTakeReward struct {
	_ struct{} `cmd_ref:"CMD_RES_QUEST_REWARD" type:"res"`
	// 今日活跃度
	DailyScore int32 `json:"dailyScore"`
	// 本周活跃度
	WeeklyScore int32 `json:"weeklyScore"`
	// 奖励vo
	RewardVos []*RewardVo `json:"rewardVos"`
}
