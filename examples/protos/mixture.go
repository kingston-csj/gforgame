package protos

type PushIdleInfo struct { // 挂机——推送挂机信息
	_         struct{} `cmd_ref:"CmdPushIdleInfo" type:"push"`
	BeginTime int64    `json:"beginTime"` // 开始挂机时间
}

type ReqClientUploadEvent struct { // 客户端上传事件
	_    struct{} `cmd_ref:"CmdReqClientUploadEvent" type:"req"`
	Type int32    `json:"type"` // 客户端自定义事件类型
}

type ResClientUploadEvent struct { // 客户端上传事件响应
	_    struct{} `cmd_ref:"CmdResClientUploadEvent" type:"res"`
	Code int32    `json:"code"` // 错误码
}

type ReqIdleGetReward struct { // 挂机——获取奖励
	_ struct{} `cmd_ref:"CmdReqIdleGetReward" type:"req"`
}

type ResIdleGetReward struct { // 挂机——获取奖励响应
	_ struct{} `cmd_ref:"CmdResIdleGetReward" type:"res"`
}

type ReqIdleViewReward struct { // 挂机——查看奖励
	_ struct{} `cmd_ref:"CmdReqIdleViewReward" type:"req"`
}

type ResIdleViewReward struct { // 挂机——查看奖励响应
	_    struct{} `cmd_ref:"CmdIdleViewReward" type:"res"`
	Code int32    `json:"code"` // 错误码
}
