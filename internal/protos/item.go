package protos

// 道具——推送背包信息
type PushBackpackInfo struct {
	_     struct{}   `cmd_ref:"CmdPushBackpackInfo"`
	Items []ItemInfo `json:"items"`
}

// 道具信息
type ItemInfo struct {
	Cf_id int32  `json:"cf_id"`
	Uid   string `json:"uid"`
	Count int32  `json:"count"`
	Level int32  `json:"level"`
	Extra string `json:"extra"`
}

// 道具——推送钱包信息
type PushPurseInfo struct {
	_       struct{} `cmd_ref:"CmdItemResPurseInfo"`
	Diamond int32    `json:"diamond"`
	Gold    int32    `json:"gold"`
}

// 道具——推送道具变化信息
type PushItemChanged struct {
	_ struct{} `cmd_ref:"CmdItemPushChanged"`
	// item, rune,card 等道具类型
	Type    string     `json:"type"`
	Changed []ItemInfo `json:"changed"`
}
