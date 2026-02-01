package protos

type ReqPlayerLogin struct {
    _        struct{} `cmd_ref:"CmdPlayerReqLogin"`
    Pwd      string
    PlayerId string
}

type ReqPlayerLoadingFinish struct{
    _        struct{} `cmd_ref:"CmdPlayerReqLoadingFinish"`
}

type ResPlayerLogin struct { // 玩家登录响应
    _          struct{} `cmd_ref:"CmdPlayerResLogin"`
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
    _     struct{} `cmd_ref:"CmdPlayerPushLoadComplete"`
}

type PushDailyResetInfo struct { // 玩家每日重置信息推送
    _     struct{} `cmd_ref:"CmdPlayerPushDailyResetInfo"`
    NormalRecruitTimes int32 `json:"normalRecruitTimes"` // 普通招募次数
    HighRecruitTimes int32 `json:"highRecruitTimes"` // 高级招募次数
    MallDailyBuyTimes int32 `json:"mallDailyBuyTimes"` // 商城每日购买次数
    DailyRechargeSum int32 `json:"dailyRechargeSum"` // 每日充值金额
}

type ReqPlayerRefreshScore struct { // 上报经营评分
    _     struct{} `cmd_ref:"CmdPlayerReqRefreshScore"`
    Score int32 `json:"score"`
}

type ResPlayerRefreshScore struct { // 玩家经营评分刷新响应
    _    struct{} `cmd_ref:"CmdPlayerResRefreshScore"`
    Code int32 `json:"code"`
}