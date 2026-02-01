package protos

type PushRechargePay struct {
	_          struct{}    `cmd_ref:"CmdPushRechargePay" type:"push"`
	RechargeId int32       `json:"rechargeId"`
	Rewards    []*RewardVo `json:"rewards"`
}

type PushRechargeInfo struct { // 充值总览
	_             struct{} `cmd_ref:"CmdPushRechargePayInfo" type:"push"`
	RechargeTimes string   `json:"rechargeTimes"` // 历史充值 id1=count1,id2=count2
}