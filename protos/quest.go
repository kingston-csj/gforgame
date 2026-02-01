package protos

type QuestVo struct {
	Id       int32 `json:"id"`
	Progress int32 `json:"progress"` // 进度
	Target   int32 `json:"target"`   // 目标
	Status   int8  `json:"status"`   // 状态 0未完成，1已完成未领奖,2已领奖
}

type PushQuestAutoTakeReward struct { // 主线任务——自动领奖
	_         struct{}    `cmd_ref:"CMD_PUSH_QUEST_AUTO_REWARD" type:"push"`
	RewardVos []*RewardVo `json:"rewardVos"`
}

type PushQuestDailyInfo struct {
	_                struct{}   `cmd_ref:"CMD_PUSH_DAILY_QUEST" type:"push"`
	DailyRewardIndex int32      `json:"dailyRewardIndex"` // 已领取的档位索引（0为未领取)
	DailyScore       int32      `json:"dailyScore"`       // 今日活跃度
	Quests           []*QuestVo `json:"quests"`           // 所有任务
}

type PushQuestRefreshVo struct {
	_     struct{} `cmd_ref:"CMD_PUSH_UPDATE_QUEST" type:"push"`
	Quest *QuestVo `json:"quest"`
}

type PushQuestReplace struct {
	_          struct{} `cmd_ref:"CMD_RES_REPLACE_QUEST" type:"push"`
	OldQuestId int32    `json:"oldQuestId"` // 旧任务id

	Quest *QuestVo `json:"quest"`
}

type PushQuestWeeklyInfo struct { // 每周任务主界面信息
	_                 struct{}   `cmd_ref:"CMD_PUSH_WEEKLY_QUEST" type:"push"`
	WeeklyRewardIndex int32      `json:"weeklyRewardIndex"` // 已领取的档位索引（0为未领取)
	WeeklyScore       int32      `json:"weeklyScore"`       // 本周活跃度
	Quests            []*QuestVo `json:"quests"`            // 所有任务
}

type PushAchievementInfo struct { // 成就——加载所有信息
	_              struct{}   `cmd_ref:"CMD_PUSH_ACHIEVEMENT" type:"push"`
	Score          int32      `json:"score"`          // 积分
	AchievementVos []*QuestVo `json:"achievementVos"` // 所有任务
}

type ReqQuestTakeAllRewards struct { // 任务——一键领取所有奖励
	_ struct{} `cmd_ref:"CMD_REQ_QUEST_ALL_REWARD" type:"req"`

	Category int32 `json:"category"` // 任务类型 1主线，2日常
}

type ReqQuestTakeProgressReward struct { // 任务——领取档位奖
	_ struct{} `cmd_ref:"CMD_REQ_QUEST_PROGRESS_REWARD" type:"req"`

	Type int32 `json:"type"` // 任务类型 2每日，5每周，6公会
}

type ReqQuestTakeReward struct { // 任务——领取达标奖
	_ struct{} `cmd_ref:"CMD_REQ_QUEST_REWARD" type:"req"`

	Id int32 `json:"id"` // 任务id
}

type ResQuestTakeAllRewards struct { // 任务——一键领取所有奖励
	_ struct{} `cmd_ref:"CMD_RES_QUEST_ALL_REWARD" type:"res"`

	Category    int32       `json:"category"`    // 任务类型 1主线，2日常
	RewardVos   []*RewardVo `json:"rewardVos"`   // 奖励vo
	DailyScore  int32       `json:"dailyScore"`  // 今日活跃度
	WeeklyScore int32       `json:"weeklyScore"` // 本周活跃度
	QuestIds    []int32     `json:"questIds"`    // 已领取的任务id列表
	Score       int32       `json:"score"`       // 总分数
}

type ResQuestTakeProgressReward struct { // 任务——领取档位奖
	_           struct{}    `cmd_ref:"CMD_RES_QUEST_PROGRESS_REWARD" type:"res"`
	Type        int32       `json:"type"`        // 任务类型 2每日，5每周，6公会
	RewardIndex int32       `json:"rewardIndex"` // 已领取的档位索引（0为未领取)
	RewardVos   []*RewardVo `json:"rewardVos"`   // 奖励vo
}

type ResQuestTakeReward struct { // 任务——领取达标奖
	_           struct{}    `cmd_ref:"CMD_RES_QUEST_REWARD" type:"res"`
	Code        int32       `json:"code"`
	DailyScore  int32       `json:"dailyScore"`  // 今日活跃度
	WeeklyScore int32       `json:"weeklyScore"` // 本周活跃度
	RewardVos   []*RewardVo `json:"rewardVos"`   // 奖励vo
}
