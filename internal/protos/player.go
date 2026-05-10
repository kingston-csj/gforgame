package protos

type ReqPlayerLogin struct {
	_        struct{} `cmd_ref:"CmdPlayerReqLogin"`
	Pwd      string   `validate:"required" json:"pwd"`
	PlayerId string   `validate:"required" json:"playerId"`
	ServerId int32    `json:"serverId"`
}

type ReqPlayerLoadingFinish struct {
	_ struct{} `cmd_ref:"CmdPlayerReqLoadingFinish"`
}

type ResPlayerLogin struct { // 玩家登录响应
	_          struct{} `cmd_ref:"CmdPlayerResLogin"`
	Code       int32    `json:"code"`
	PlayerId   string   `json:"playerId"`
	NewCreate  bool     `json:"newCreate"`
	CreateTime int64    `json:"createTime"`
	Head       int32    `json:"head"`
	Level      int32    `json:"level"`
	Stage      int32    `json:"stage"`
	Name       string   `json:"name"`
	Fighting   int32    `json:"fighting"`
	Camp       int32    `json:"camp"`
}

type ReqPlayerCreate struct {
	_    struct{} `cmd_ref:"CmdPlayerReqCreate" type:"req"`
	Name string   `json:"name"`
	Camp int32    `json:"camp"`
}

type ResPlayerCreate struct {
	_        struct{} `cmd_ref:"CmdPlayerResCreate" type:"res"`
	Code     int32    `json:"code"`
	PlayerId string   `json:"playerId"`
}

type PushReplacingLogin struct {
}

type ReqJoinRoom struct {
	_        struct{} `cmd_ref:"CmdChaJoinRoom" type:"req"`
	RoomId   int64    `json:"roomId"`
	PlayerId int64    `json:"playerId"`
}

type ReqPlayerUpLevel struct {
	_       struct{} `cmd_ref:"CmdPlayerReqUpLevel" type:"req"`
	ToLevel int32    `json:"toLevel"`
}

type ResPlayerUpLevel struct {
	_    struct{} `cmd_ref:"CmdPlayerResUpLevel" type:"res"`
	Code int32    `json:"code"`
}

type PushPlayerFightChange struct {
	_     struct{} `cmd_ref:"CmdPlayerPushFightChange" type:"push"`
	Fight int32    `json:"fight"`
}

type ReqPlayerUpStage struct {
	_ struct{} `cmd_ref:"CmdPlayerReqUpStage" type:"req"`
}

type ResPlayerUpStage struct {
	_    struct{} `cmd_ref:"CmdPlayerResUpStage" type:"res"`
	Code int32    `json:"code"`
}

type PushLoadComplete struct {
	_ struct{} `cmd_ref:"CmdPlayerPushLoadComplete"`
}

type PushDailyResetInfo struct { // 玩家每日重置信息推送
	_                  struct{} `cmd_ref:"CmdPlayerPushDailyResetInfo"`
	NormalRecruitTimes int32    `json:"normalRecruitTimes"` // 普通招募次数
	HighRecruitTimes   int32    `json:"highRecruitTimes"`   // 高级招募次数
	MallDailyBuyTimes  int32    `json:"mallDailyBuyTimes"`  // 商城每日购买次数
	DailyRechargeSum   int32    `json:"dailyRechargeSum"`   // 每日充值金额
}

type PushWeeklyResetInfo struct { // 玩家每周重置信息推送
	_ struct{} `cmd_ref:"CmdPushWeeklyResetInfo"`

	WeeklyRechargeSum float32 `json:"weeklyRechargeSum"` // 充值累计积分
	WeeklyGiftRewards string  `json:"weeklyGiftRewards"` // 充值礼包奖励领取状态, 格式为 id1=status,id2=status
	RechargeBuyTimes  string  `json:"rechargeBuyTimes"`  // 充值每周限购商品信息, 格式为 id1=次数1,id2=次数2
}

type ReqPlayerRefreshScore struct { // 上报经营评分
	_     struct{} `cmd_ref:"CmdPlayerReqRefreshScore"`
	Score int32    `json:"score"`
}

type ResPlayerRefreshScore struct { // 玩家经营评分刷新响应
	_    struct{} `cmd_ref:"CmdPlayerResRefreshScore"`
	Code int32    `json:"code"`
}

type ReqEditClientData struct { // 玩家编辑客户端数据请求
	_    struct{} `cmd_ref:"CmdPlayerReqEditClientData"`
	Data string   `json:"data"`
}

type ResEditClientData struct { // 玩家编辑客户端数据响应
	_    struct{} `cmd_ref:"CmdPlayerResEditClientData"`
	Code int32    `json:"code"`
}