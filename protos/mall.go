package protos

// 商城——主页信息
type PushMallInfo struct {
	_     struct{} `cmd_ref:"CmdPushMallInfo" type:"push"`
	LiftBuyTimes map[int32]int32 `json:"liftBuyTimes"` // 商城终身限购次数
}

// 商城——购买
type ReqMallBuy struct {
	_     struct{} `cmd_ref:"CmdReqMallBuy" type:"req"`
	ProductId int32 `json:"productId"` // 物品ID
	Count  int32 `json:"count"`  // 购买数量
}

// 商城——购买结果
type ResMallBuy struct {
	_     struct{} `cmd_ref:"CmdResMallBuy" type:"res"`
	Code int32 `json:"code"` // 错误码
}