package protos

type ReqHeartBeat struct { // 心跳请求
	_     struct{} `cmd_ref:"CmdHeartBeatReq"`
	Index int32    `json:"index"`
}

type ResHeartBeat struct { // 心跳请求
	_     struct{} `cmd_ref:"CmdHeartBeatRes"`
	Index int32    `json:"index"`
	Code  int      `json:"code"`
}

type ReqGetServerTime struct { // 获取服务器时间
	_     struct{} `cmd_ref:"CmdGetServerTimeReq"`
	Index int32    `json:"index"`
}

type ResGetServerTime struct { // 获取服务器时间
	_          struct{} `cmd_ref:"CmdGetServerTimeRes"`
	ServerTime int64    `json:"serverTime"`
	Code       int      `json:"code"`
}