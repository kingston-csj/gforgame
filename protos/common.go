package protos

// 奖励vo
// Item id_count;
// Hero id
type RewardVo struct {
	Type  int32 `json:"type"`
	Value int32 `json:"value"`
}