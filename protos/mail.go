package protos

type MailVo struct {
	Id int64 `json:"id"`
	// 邮件标题， 当TemplateId为0时，需要此字段
	Title string `json:"title"`
	// 邮件内容， 当TemplateId为0时，需要此字段
	Content string `json:"content"`
	// 邮件奖励
	Rewards []RewardInfo `json:"rewards"`
	// 邮件模板id
	TemplateId int32 `json:"templateId"`
	// 邮件状态
	Status int32 `json:"status"`
	// 邮件时间
	Time int64 `json:"time"`
}

type PushMailAll struct {
	Mails []MailVo `json:"mails"`
}

type ReqMailGetAllRewards struct{}

type ResMailGetAllRewards struct {
	Code int32 `json:"code"`
}

type ReqMailRead struct {
	Id int64 `json:"id"`
}

type ResMailRead struct {
	Code int32 `json:"code"`
}

type ReqMailGetReward struct {
	Id int64 `json:"id"`
}

type ResMailGetReward struct {
	Code int32 `json:"code"`
}

type ReqMailDeleteAll struct{}

type ResMailDeleteAll struct {
	Removed []int64 `json:"removed"`
}

type ResMailDelete struct {
	Code int32 `json:"code"`
}

type ReqMailReadAll struct{}

type ResMailReadAll struct {
	Code int32 `json:"code"`
}
