package player

type FriendApplyItem struct {
	FromId   string `json:"player_id"`
	FromName string `json:"player_name"`
	TargetId string `json:"target_id"`
	Status   int32  `json:"status"`
	Id       string `json:"id"`
	Time     int64  `json:"time"`
	Fighting int64  `json:"fighting"`
}
