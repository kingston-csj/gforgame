package protos

type FriendVo struct {
	ApplyId  string `json:"applyId"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Head     int    `json:"head"`
	Fighting int64  `json:"fighting"`
}

type FriendApplyVo struct {
	Id         string `json:"id"`
	FromId     string `json:"fromId"`
	FromName   string `json:"fromName"`
	TargetId   string `json:"targetId"`
	TargetName string `json:"targetName"`
	Status int   `json:"status"` // 申请结果：1同意 2拒绝 0未处理
	Time   int64 `json:"time"`
}

type PushFriendInfo struct {
    _           struct{} `cmd_ref:"CmdFriendPushApplyList"`
    ApplyItems  []*FriendApplyVo `json:"applyItems"`
    FriendItems []*FriendVo      `json:"friendItems"`
    FriendSum   int              `json:"friendSum"`
}

type ReqFriendSearchPlayers struct {
    _   struct{} `cmd_ref:"CmdFriendReqSearchPlayers"`
    Key string `json:"key"`
}

type ReqFriendApply struct {
    _        struct{} `cmd_ref:"CmdFriendReqApply"`
    FriendId string `json:"friendId"`
}

type ReqFriendDealApplyRecord struct {
    _       struct{} `cmd_ref:"CmdFriendReqDealApply"`
    ApplyId string `json:"applyId"`
    Status  int    `json:"status"`
}

type ReqFriendDelete struct {
    _        struct{} `cmd_ref:"CmdFriendReqDelete"`
    FriendId string `json:"friendId"`
}

type ReqFriendQueryMyFriends struct {
    _ struct{} `cmd_ref:"CmdFriendReqQueryFriends"`
}

type ResFriendApply struct {
    _    struct{} `cmd_ref:"CmdFriendResApply"`
    Code int `json:"code"`
}

type ResFriendDealApplyRecord struct {
    _    struct{} `cmd_ref:"CmdFriendResDealApply"`
    Code int `json:"code"`
}

type ResFriendDelete struct {
    _    struct{} `cmd_ref:"CmdFriendResDelete"`
    Code int `json:"code"`
}

// 查询我的好友列表
type ResFriendQueryMyFriends struct {
    _     struct{} `cmd_ref:"CmdFriendResQueryFriends"`
    Code  int         `json:"code"`
    Items []*FriendVo `json:"items"`
}

// 通过id或昵称搜索玩家
type ResFriendSearchPlayers struct {
    _     struct{} `cmd_ref:"CmdFriendResSearchPlayers"`
    Code  int         `json:"code"`
    Items []*FriendVo `json:"items"`
}
