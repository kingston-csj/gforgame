package protos

type ReqPlayerLogin struct {
	Id  string
	Pwd string
}

type ReqPlayerLoadingFinish struct{}

type ResPlayerLogin struct {
	Code     int32  `json:"code"`
	Name     string `json:"name"`
	Fighting int32  `json:"fighting"`
	Camp     int32  `json:"camp"`
}

type ReqPlayerCreate struct {
	Name string
}

type ResPlayerCreate struct {
	Id int64
}

type ReqJoinRoom struct {
	RoomId int64

	PlayerId int64
}

type ReqChat struct {
	Id string
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
