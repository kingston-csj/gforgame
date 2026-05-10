package protos

type PushMonthCardInfo struct { // 月卡——主界面信息
	_          struct{}       `cmd_ref:"CmdPushMonthCardInfo"`
	SilverCard *MonthlyCardVo `json:"silverCard"`
	GoldCard   *MonthlyCardVo `json:"goldCard"`
}

type MonthlyCardVo struct {
	ExpiredTime int64 `json:"expiredTime"` // 月卡过期时间
}

type ReqMonthCardGetReward struct {
	_    struct{} `cmd_ref:"CmdReqMonthCardGetReward"`
	Type int32    `json:"type"` // 月卡类型
}

type ResMonthCardGetReward struct {
	_    struct{} `cmd_ref:"CmdResMonthCardGetReward"`
	Code int32    `json:"code"`
}
