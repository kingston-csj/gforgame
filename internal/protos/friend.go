package protos

type FriendVo struct { // 好友信息
	ApplyId  string `json:"applyId"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Head     int32  `json:"head"`
	Fighting int64  `json:"fighting"`
}

type FriendApplyVo struct { // 好友申请信息
	Id         string `json:"id"`
	FromId     string `json:"fromId"`
	FromName   string `json:"fromName"`
	FromHead   int32  `json:"fromHead"`
	TargetId   string `json:"targetId"`
	TargetName string `json:"targetName"`
	Status     int32  `json:"status"` // 申请结果：1同意 2拒绝 0未处理
	Time       int64  `json:"time"`
}

type PushFriendInfo struct { // 好友——加载所有数据
	_           struct{}         `cmd_ref:"CmdFriendPushApplyList"`
	ApplyItems  []*FriendApplyVo `json:"applyItems"`
	FriendItems []*FriendVo      `json:"friendItems"`
	FriendSum   int32            `json:"friendSum"`
}

type ReqFriendSearchPlayers struct { // 好友——搜索玩家
	_   struct{} `cmd_ref:"CmdFriendReqSearchPlayers"`
	Key string   `json:"key"` // 搜索关键词
}

type ReqFriendApply struct { // 好友——申请好友
	_        struct{} `cmd_ref:"CmdFriendReqApply"`
	FriendId string   `json:"friendId"`
}

type ReqFriendDealApplyRecord struct { // 好友——处理申请
	_       struct{} `cmd_ref:"CmdFriendReqDealApply"`
	ApplyId string   `json:"applyId"`
	Status  int32    `json:"status"` // 处理结果：1同意 2拒绝
}

type ReqFriendDelete struct { // 好友——删除
	_        struct{} `cmd_ref:"CmdFriendReqDelete"`
	FriendId string   `json:"friendId"` // 好友id
}

type ReqFriendQueryMyFriends struct { // 好友——查询我的好友
	_ struct{} `cmd_ref:"CmdFriendReqQueryFriends"`
}

type ResFriendApply struct { // 好友——申请好友
	_    struct{} `cmd_ref:"CmdFriendResApply"`
	Code int32    `json:"code"`
}

type ResFriendDealApplyRecord struct { // 好友——处理申请
	_    struct{} `cmd_ref:"CmdFriendResDealApply"`
	Code int32    `json:"code"`
}

type ResFriendDelete struct { // 好友——删除
	_    struct{} `cmd_ref:"CmdFriendResDelete"`
	Code int32    `json:"code"`
}

// 查询我的好友列表
type ResFriendQueryMyFriends struct { // 好友——查询我的好友
	_     struct{}    `cmd_ref:"CmdFriendResQueryFriends"`
	Code  int         `json:"code"`
	Items []*FriendVo `json:"items"`
}

// 通过id或昵称搜索玩家
type ResFriendSearchPlayers struct {
	_     struct{}    `cmd_ref:"CmdFriendResSearchPlayers"`
	Code  int         `json:"code"`
	Items []*FriendVo `json:"items"` // 搜索结果
}
