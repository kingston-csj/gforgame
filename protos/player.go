package protos

type ReqPlayerLogin struct {
    _        struct{} `cmd_ref:"CmdPlayerReqLogin" type:"req"`
    Pwd      string
    PlayerId string
}

type ReqPlayerLoadingFinish struct{
    _        struct{} `cmd_ref:"CmdPlayerReqLoadingFinish" type:"req"`
}

// ResPlayerLogin 玩家登录响应
type ResPlayerLogin struct {
    _          struct{} `cmd_ref:"CmdPlayerResLogin" type:"res"`
    Code       int32  `json:"code"`
    PlayerId   string `json:"playerId"`
    NewCreate  bool   `json:"newCreate"`
    CreateTime int64  `json:"createTime"`
    Head       int32  `json:"head"`
    Level      int32  `json:"level"`
    Name       string `json:"name"`
    Fighting   int32  `json:"fighting"`
    Camp       int32  `json:"camp"`
}

type ReqPlayerCreate struct {
    _    struct{} `cmd_ref:"CmdPlayerReqCreate" type:"req"`
    Name string `json:"name"`
    Camp int32 `json:"camp"`
}

type ResPlayerCreate struct {
    _        struct{} `cmd_ref:"CmdPlayerResCreate" type:"res"`
    Code     int32 `json:"code"`
    PlayerId string `json:"playerId"`
}

type PushReplacingLogin struct {
}

type ReqJoinRoom struct {
    _        struct{} `cmd_ref:"CmdChaJoinRoom" type:"req"`
    RoomId   int64 `json:"roomId"`
    PlayerId int64 `json:"playerId"`
}

type ReqPlayerUpLevel struct {
    _       struct{} `cmd_ref:"CmdPlayerReqUpLevel" type:"req"`
    ToLevel int32 `json:"toLevel"`
}

type ResPlayerUpLevel struct {
    _    struct{} `cmd_ref:"CmdPlayerResUpLevel" type:"res"`
    Code int32 `json:"code"`
}

type PushPlayerFightChange struct {
    _     struct{} `cmd_ref:"CmdPlayerPushFightChange" type:"push"`
    Fight int32 `json:"fight"`
}

type ReqPlayerUpStage struct{
    _     struct{} `cmd_ref:"CmdPlayerReqUpStage" type:"req"`
}

type ResPlayerUpStage struct {
    _    struct{} `cmd_ref:"CmdPlayerResUpStage" type:"res"`
    Code int32 `json:"code"`
}

type PushLoadComplete struct {
    _     struct{} `cmd_ref:"CmdPlayerPushLoadComplete" type:"push"`
}