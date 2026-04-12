package protos

// 签到——加载所有信息
type PushSigninInfo struct {
	_ struct{} `cmd_ref:"CmdSignInPush" type:"push"`
	DaysInMonth int32 `json:"daysInMonth"` // 本月总天数
	NthDay      int32 `json:"nthDay"`      // 今天第几天
	SigninDays  []int32 `json:"signinDays"` // 已签到天数
	SignInMakeUp map[int32]int32 `json:"signInMakeUp"` // 是否已补签
}

// 签到——请求
type ReqSignIn struct {
	_ struct{} `cmd_ref:"CmdSignInReqSignIn"`
}

// 签到——响应
type ResSignIn struct {
	_ struct{} `cmd_ref:"CmdSignInResSignIn"`
	Code int32 `json:"code"`
}

type ReqSignInMakeup struct { // 签到--补签
	_ struct{} `cmd_ref:"CmdSignInReqSignInMakeup"`	
	Day int32 `json:"day"` // 补签的目标天
}

type ResSignInMakeup struct {// 签到--补签
	_ struct{} `cmd_ref:"CmdSignInResSignInMakeup"`
	Code int32 `json:"code"`
}