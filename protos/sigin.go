package protos

// 签到——加载所有信息
type PushSigninInfo struct {
	_ struct{} `cmd_ref:"CmdSignInPush" type:"push"`
	DaysInMonth int32 `json:"daysInMonth"` // 本月总天数
	NthDay      int32 `json:"nthDay"`      // 今天第几天
	SigninDays  []int32 `json:"signinDays"` // 已签到天数
}

// 签到——请求
type ReqSignIn struct {
	_ struct{} `cmd_ref:"CmdSignInReqSignIn" type:"req"`
}

// 签到——响应
type ResSignIn struct {
	_ struct{} `cmd_ref:"CmdSignInResSignIn" type:"res"`
	Code int32 `json:"code"`
}