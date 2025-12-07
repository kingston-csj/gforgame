// nothging
package protos

type ChatMessageVo struct  // 聊天消息vo
{ 
	Id         string `json:"id"` // 消息id
	Channel    int32  `json:"channel"` // 发送频道：1个人 2世界
	SenderId   string `json:"senderId"`
	SenderName string `json:"senderName"`
	SenderHead int    `json:"senderHead"`
	ReceiverId string `json:"receiverId"`
	Timestamp  int64  `json:"timestamp"`
	Content    string `json:"content"`
}

// / PushChatNewMessage 推送新聊天消息
type PushChatNewMessage struct {
    _        struct{} `cmd_ref:"CmdChatPushNew" type:"push"`
    Code     int              `json:"code"`
    Messages []*ChatMessageVo `json:"messages"`
}

// / ReqChat 发送聊天消息请求
type ReqChat struct {
    _       struct{} `cmd_ref:"CmdChatReqChat" type:"req"`
    Channel int    `json:"channel"`
    Target  string `json:"target"`
    Content string `json:"content"`
}

type ResChat struct {
    _    struct{} `cmd_ref:"CmdChatResChat" type:"res"`
    Code int `json:"code"`
}
