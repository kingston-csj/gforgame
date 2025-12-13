// nothging
package protos

type MailVo struct { // 邮件vo
	Id         int64        `json:"id"`
	Title      string       `json:"title"`      // 邮件标题， 当TemplateId为0时，需要此字段
	Content    string       `json:"content"`    // 邮件内容， 当TemplateId为0时，需要此字段
	Rewards    []RewardVo `json:"rewards"`    // 邮件奖励
	TemplateId int32        `json:"templateId"` // 邮件模板id
	Status     int32        `json:"status"`     // 邮件状态
	Time       int64        `json:"time"`       // 邮件时间
}

type PushMailAll struct {
    _     struct{} `cmd_ref:"CmdMailPushAll" type:"push"`
    Mails []MailVo `json:"mails"`
}

type ReqMailGetAllRewards struct{
    _ struct{} `cmd_ref:"CmdMailReqGetAllReward" type:"req"`
}

type ResMailGetAllRewards struct {
    _    struct{} `cmd_ref:"CmdMailResGetAllReward" type:"res"`
    Code int32 `json:"code"`
}

type ReqMailRead struct {
    _  struct{} `cmd_ref:"CmdMailReqRead" type:"req"`
    Id int64 `json:"id"`
}

type ResMailRead struct {
    _    struct{} `cmd_ref:"CmdMailResRead" type:"res"`
    Code int32 `json:"code"`
}

type ReqMailGetReward struct {
    _  struct{} `cmd_ref:"CmdMailReqGetReward" type:"req"`
    Id int64 `json:"id"`
}

type ResMailGetReward struct {
    _    struct{} `cmd_ref:"CmdMailResGetReward" type:"res"`
    Code int32 `json:"code"`
}

type ReqMailDeleteAll struct{
    _ struct{} `cmd_ref:"CmdMailReqDeleteAll" type:"req"`
}

type ResMailDeleteAll struct {
    _       struct{} `cmd_ref:"CmdMailResDeleteAll" type:"res"`
    Removed []int64 `json:"removed"`
}

type ResMailDelete struct {
	Code int32 `json:"code"`
}

type ReqMailReadAll struct{
    _ struct{} `cmd_ref:"CmdMailReqReadAll" type:"req"`
}

type ResMailReadAll struct {
    _    struct{} `cmd_ref:"CmdMailResReadAll" type:"res"`
    Code int32 `json:"code"`
}
