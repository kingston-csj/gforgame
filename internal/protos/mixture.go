package protos

type ReqClientUploadEvent struct { // 客户端上传事件
	_    struct{} `cmd_ref:"CmdReqClientUploadEvent" type:"req"`
	Type int32    `json:"type"` // 客户端自定义事件类型
}

type ResClientUploadEvent struct { // 客户端上传事件响应
	_    struct{} `cmd_ref:"CmdResClientUploadEvent" type:"res"`
	Code int32    `json:"code"` // 错误码
}
