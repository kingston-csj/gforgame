package protos

// 月卡——主界面信息
type PushMonthCardInfo struct {
	_     struct{} `cmd_ref:"CmdPushMonthCardInfo" type:"push"`
	SilverCard *MonthlyCardVo `json:"silverCard"`
	GoldCard *MonthlyCardVo `json:"goldCard"`
}

type MonthlyCardVo struct {
	ExpiredTime int64 `json:"expiredTime"` // 月卡过期时间
}

type ReqMonthCardGetReward struct {
	_ struct{} `cmd_ref:"CmdReqMonthCardGetReward" type:"req"`
	Type int32 `json:"type"` // 月卡类型
}

type ResMonthCardGetReward struct {
	_ struct{} `cmd_ref:"CmdResMonthCardGetReward" type:"res"`
	Code int32 `json:"code"`
}
