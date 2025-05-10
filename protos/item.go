package protos

type ResBackpackInfo struct {
	Items []ItemInfo `json:"items"`
}

type ItemInfo struct {
	Id    int32 `json:"id"`
	Count int32 `json:"count"`
}

type ReqGmAction struct {
	Topic  string
	Params string
}

type ResGmAction struct {
	Code int32 `json:"code"`
}

type RewardInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PushPurseInfo struct {
	Diamond int32 `json:"diamond"`
	Gold    int32 `json:"gold"`
}

type PushItemChanged struct {
	ItemId int32 `json:"itemId"`
	Count  int32 `json:"count"`
}
