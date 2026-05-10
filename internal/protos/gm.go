package protos

type ReqGmCommand struct { //设置gm命令
	_    struct{} `cmd_ref:"CmdGmReqCommand"`
	Args string   `json:"args"` //添加道具 add_items 1001=1;1=2添加货币 add_currency Diamond_100设置玩家等级 level 50
}

type ResGmCommand struct { //设置gm命令
	_    struct{} `cmd_ref:"CmdGmResCommand"`
	Code int32    `json:"code"`
}
