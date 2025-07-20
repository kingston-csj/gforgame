package protos

type ChatMessageVo struct {
	Id         string `json:"id"`
	Channel    int32  `json:"channel"`
	SenderId   string `json:"senderId"`
	SenderName string `json:"senderName"`
	SenderHead int    `json:"senderHead"`
	ReceiverId string `json:"receiverId"`
	Timestamp  int64  `json:"timestamp"`
	Content    string `json:"content"`
}

type PushChatNewMessage struct {
	Code     int              `json:"code"`
	Messages []*ChatMessageVo `json:"messages"`
}

type ReqChat struct {
	// 发送频道：1个人 2世界
	Channel int `json:"channel"`
	// 发送目标：个人id 世界id
	Target string `json:"target"`
	// 发送内容
	Content string `json:"content"`
}

type ResChat struct {
	Code int `json:"code"`
}
