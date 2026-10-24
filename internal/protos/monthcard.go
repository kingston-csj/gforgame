package protos

// 月卡——主界面信息
type PushMonthCardInfo struct { 
	_          struct{}       `cmd_ref:"CmdPushMonthCardInfo"`
	SilverCard *MonthlyCardVo `json:"silverCard"`
	GoldCard   *MonthlyCardVo `json:"goldCard"`
}

// 月卡信息
type MonthlyCardVo struct {
	ExpiredTime int64 `json:"expiredTime"` // 月卡过期时间
}

// 月卡——获取奖励
type ReqMonthCardGetReward struct {
	_    struct{} `cmd_ref:"CmdReqMonthCardGetReward"`
	Type int32    `json:"type"` // 月卡类型
}

// 月卡——获取奖励响应
type ResMonthCardGetReward struct {
	_    struct{} `cmd_ref:"CmdResMonthCardGetReward"`
	Code int32    `json:"code"`
}
