package protos

type ReqPlayerLogin struct {
	Pwd      string
	PlayerId string
}

type ReqPlayerLoadingFinish struct{}

// ResPlayerLogin 玩家登录响应
type ResPlayerLogin struct {
	Code       int32  `json:"code"`
	PlayerId   string `json:"playerId"`   //  玩家ID
	NewCreate  bool   `json:"newCreate"`  //  是否是新创建的玩家
	CreateTime int64  `json:"createTime"` //  创建时间
	Head       int32  `json:"head"`       //  玩家头像
	Level      int32  `json:"level"`      //  玩家等级
	Name       string `json:"name"`       //  玩家名称
	Fighting   int32  `json:"fighting"`   //  玩家战斗力
	Camp       int32  `json:"camp"`       //  玩家阵营
}

type ReqPlayerCreate struct {
	Name string
}

type ResPlayerCreate struct {
	Id int64
}

type PushReplacingLogin struct {
}

type ReqJoinRoom struct {
	RoomId   int64 `json:"roomId"`   //  房间ID
	PlayerId int64 `json:"playerId"` //  玩家ID
}

type ReqPlayerUpLevel struct {
	ToLevel int32 `json:"toLevel"`
}

type ResPlayerUpLevel struct {
	Code int32 `json:"code"`
}

type PushPlayerFightChange struct {
	Fight int32 `json:"fight"`
}

type ReqPlayerUpStage struct{}

type ResPlayerUpStage struct {
	Code int32 `json:"code"`
}
