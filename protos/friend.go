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
	// 申请结果：1同意 2拒绝 0未处理
	Status int   `json:"status"`
	Time   int64 `json:"time"`
}

type PushFriendInfo struct {
    _           struct{} `cmd_ref:"CmdFriendPushApplyList" type:"push"`
    ApplyItems  []*FriendApplyVo `json:"applyItems"`
    FriendItems []*FriendVo      `json:"friendItems"`
    FriendSum   int              `json:"friendSum"`
}

type ReqFriendSearchPlayers struct {
    _   struct{} `cmd_ref:"CmdFriendReqSearchPlayers" type:"req"`
    Key string `json:"key"`
}

type ReqFriendApply struct {
    _        struct{} `cmd_ref:"CmdFriendReqApply" type:"req"`
    FriendId string `json:"friendId"`
}

type ReqFriendDealApplyRecord struct {
    _       struct{} `cmd_ref:"CmdFriendReqDealApply" type:"req"`
    ApplyId string `json:"applyId"`
    Status  int    `json:"status"`
}

type ReqFriendDelete struct {
    _        struct{} `cmd_ref:"CmdFriendReqDelete" type:"req"`
    FriendId string `json:"friendId"`
}

type ReqFriendQueryMyFriends struct {
    _ struct{} `cmd_ref:"CmdFriendReqQueryFriends" type:"req"`
}

type ResFriendApply struct {
    _    struct{} `cmd_ref:"CmdFriendResApply" type:"res"`
    Code int `json:"code"`
}

type ResFriendDealApplyRecord struct {
    _    struct{} `cmd_ref:"CmdFriendResDealApply" type:"res"`
    Code int `json:"code"`
}

type ResFriendDelete struct {
    _    struct{} `cmd_ref:"CmdFriendResDelete" type:"res"`
    Code int `json:"code"`
}

// 查询我的好友列表
type ResFriendQueryMyFriends struct {
    _     struct{} `cmd_ref:"CmdFriendResQueryFriends" type:"res"`
    Code  int         `json:"code"`
    Items []*FriendVo `json:"items"`
}

// 通过id或昵称搜索玩家
type ResFriendSearchPlayers struct {
    _     struct{} `cmd_ref:"CmdFriendResSearchPlayers" type:"res"`
    Code  int         `json:"code"`
    Items []*FriendVo `json:"items"`
}
