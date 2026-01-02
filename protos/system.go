package protos

type ReqHeartBeat struct { // 心跳请求
	_     struct{} `cmd_ref:"CmdHeartBeatReq" type:"req"`
	Index int32    `json:"index"`
}

type ResHeartBeat struct { // 心跳请求
	_     struct{} `cmd_ref:"CmdHeartBeatRes" type:"res"`
	Index int32    `json:"index"`
	Code  int      `json:"code"`
}

type ReqGetServerTime struct { // 获取服务器时间
	_     struct{} `cmd_ref:"CmdGetServerTimeReq" type:"req"`
	Index int32    `json:"index"`
}

type ResGetServerTime struct { // 获取服务器时间
	_          struct{} `cmd_ref:"CmdGetServerTimeRes" type:"res"`
	ServerTime int64    `json:"serverTime"`
	Code       int      `json:"code"`
}