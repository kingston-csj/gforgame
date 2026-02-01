package protos

type PushBackpackInfo struct {
    _     struct{} `cmd_ref:"CmdItemPushBackpackInfo"`
    Items []ItemInfo `json:"items"`
}

type ItemInfo struct {
	Cf_id    int32 `json:"cf_id"`
    Uid      string `json:"uid"`
	Count int32 `json:"count"`
    Level int32 `json:"level"`
    Extra string `json:"extra"`
}

type PushPurseInfo struct {
    _       struct{} `cmd_ref:"CmdItemResPurseInfo"`
    Diamond int32 `json:"diamond"`
    Gold    int32 `json:"gold"`
}

type PushItemChanged struct {
    _      struct{} `cmd_ref:"CmdItemPushChanged"`
    // item, rune,card 等道具类型
    Type string `json:"type"`   
    Changed []ItemInfo `json:"changed"`
}
