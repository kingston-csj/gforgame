package protos

type MailVo struct { // 邮件vo
	Id         string        `json:"id"`
	Title      string       `json:"title"`      // 邮件标题， 当TemplateId为0时，需要此字段
	Content    string       `json:"content"`    // 邮件内容， 当TemplateId为0时，需要此字段
	Rewards    []RewardVo `json:"rewards"`      // 邮件奖励
	TemplateId int32        `json:"templateId"` // 邮件模板id
	Status     int32        `json:"status"`     // 邮件状态
	Time       int64        `json:"time"`       // 邮件时间
}

type PushMailAll struct {
    _     struct{} `cmd_ref:"CmdMailPushAll"`
    Mails []MailVo `json:"mails"`
}

type ReqMailGetAllRewards struct{
    _ struct{} `cmd_ref:"CmdMailReqGetAllReward"`
}

type ResMailGetAllRewards struct {
    _    struct{} `cmd_ref:"CmdMailResGetAllReward"`
    Code int32 `json:"code"`
    Rewards []*RewardVo `json:"rewards"`
}

type ReqMailRead struct {
    _  struct{} `cmd_ref:"CmdMailReqRead"`
    Id string `json:"id"`
}

type ResMailRead struct {
    _    struct{} `cmd_ref:"CmdMailResRead"`    
    Code int32 `json:"code"`
}

type ReqMailGetReward struct {
    _  struct{} `cmd_ref:"CmdMailReqGetReward"`
    Id string `json:"id"`
}

type ResMailGetReward struct {
    _    struct{} `cmd_ref:"CmdMailResGetReward"`    
    Code int32 `json:"code"`
    Rewards []*RewardVo `json:"rewards"`
}

type ReqMailDeleteAll struct{
    _ struct{} `cmd_ref:"CmdMailReqDeleteAll"`
}

type ResMailDeleteAll struct {
    _       struct{} `cmd_ref:"CmdMailResDeleteAll"`
    Removed []string `json:"removed"`
}

type ResMailDelete struct {
	Code int32 `json:"code"`
}

type ReqMailReadAll struct{
    _ struct{} `cmd_ref:"CmdMailReqReadAll"`
}

type ResMailReadAll struct {
    _    struct{} `cmd_ref:"CmdMailResReadAll"`    
    Code int32 `json:"code"`
}
