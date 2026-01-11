package protos

// ReqGmCommand 设置gm命令
type ReqGmCommand struct {
	_    struct{} `cmd_ref:"CmdGmReqCommand" type:"req"`
	Args string   `json:"args"` //添加道具 add_items 1001=1;1=2添加货币 add_currency Diamond_100设置玩家等级 level 50
}

// ResGmCommand 设置gm命令
type ResGmCommand struct {
	_    struct{} `cmd_ref:"CmdGmResCommand" type:"res"`
	Code int32    `json:"code"`
}
