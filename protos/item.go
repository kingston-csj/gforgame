package protos

type ResBackpackInfo struct {
    _     struct{} `cmd_ref:"CmdItemResBackpackInfo" type:"res"`
    Items []ItemInfo `json:"items"`
}

type ItemInfo struct {
	Id    int32 `json:"id"`
	Count int32 `json:"count"`
}

type ReqGmAction struct {
    _      struct{} `cmd_ref:"CmdGmReqAction" type:"req"`
    Topic  string
    Params string
}

type ResGmAction struct {
    _    struct{} `cmd_ref:"CmdGmResAction" type:"res"`
    Code int32 `json:"code"`
}

type RewardInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PushPurseInfo struct {
    _       struct{} `cmd_ref:"CmdItemResPurseInfo" type:"push"`
    Diamond int32 `json:"diamond"`
    Gold    int32 `json:"gold"`
}

type PushItemChanged struct {
    _      struct{} `cmd_ref:"CmdItemPushChanged" type:"push"`
    ItemId int32 `json:"itemId"`
    Count  int32 `json:"count"`
}
