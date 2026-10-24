package protos

// 玩家登录请求
type ReqPlayerLogin struct {
	_        struct{} `cmd_ref:"CmdReqPlayerLogin"`
	Pwd      string   `validate:"required" json:"pwd"`
	PlayerId string   `validate:"required" json:"playerId"`
	ServerId int32    `json:"serverId"`
}

// 玩家加载界面完成
type ReqPlayerLoadingFinish struct {
	_ struct{} `cmd_ref:"CmdPlayerReqLoadingFinish"`
}

// 玩家登录响应
type ResPlayerLogin struct { // 玩家登录响应
	_          struct{} `cmd_ref:"CmdResPlayerLogin"`
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

// 玩家创角
type ResPlayerCreate struct {
	_        struct{} `cmd_ref:"CmdResPlayerCreate" type:"res"`
	Code     int32    `json:"code"`
	PlayerId string   `json:"playerId"`
}

// 玩家重复登录
type PushReplacingLogin struct {
}

// 玩家升级
type ReqPlayerUpLevel struct {
	_       struct{} `cmd_ref:"CmdReqPlayerUpLevel" type:"req"`
	ToLevel int32    `json:"toLevel"`
}

// 玩家升级响应
type ResPlayerUpLevel struct {
	_    struct{} `cmd_ref:"CmdResPlayerUpLevel" type:"res"`
	Code int32    `json:"code"`
}

// 玩家战斗值变化
type PushPlayerFightChange struct {
	_     struct{} `cmd_ref:"CmdPlayerPushFightChange" type:"push"`
	Fight int32    `json:"fight"`
}

// 玩家升级阶段
type ReqPlayerUpStage struct {
	_ struct{} `cmd_ref:"CmdReqPlayerUpStage" type:"req"`
}

// 玩家升级阶段响应
type ResPlayerUpStage struct {
	_    struct{} `cmd_ref:"CmdResPlayerUpStage" type:"res"`
	Code int32    `json:"code"`
}

// 玩家常规数据推送完成
type PushLoadComplete struct {
	_ struct{} `cmd_ref:"CmdPushLoadComplete"`
}

// 玩家每日重置信息推送
type PushDailyResetInfo struct { 
	_                  struct{} `cmd_ref:"CmdPushDailyResetInfo"`
	NormalRecruitTimes int32    `json:"normalRecruitTimes"` // 普通招募次数
	HighRecruitTimes   int32    `json:"highRecruitTimes"`   // 高级招募次数
	MallDailyBuyTimes  int32    `json:"mallDailyBuyTimes"`  // 商城每日购买次数
	DailyRechargeSum   int32    `json:"dailyRechargeSum"`   // 每日充值金额
}

// 玩家每周重置信息推送
type PushWeeklyResetInfo struct { 
	_ struct{} `cmd_ref:"CmdPushWeeklyResetInfo"`

	WeeklyRechargeSum float32 `json:"weeklyRechargeSum"` // 充值累计积分
	WeeklyGiftRewards string  `json:"weeklyGiftRewards"` // 充值礼包奖励领取状态, 格式为 id1=status,id2=status
	RechargeBuyTimes  string  `json:"rechargeBuyTimes"`  // 充值每周限购商品信息, 格式为 id1=次数1,id2=次数2
}

// 玩家——上报经营评分
type ReqPlayerRefreshScore struct { 
	_     struct{} `cmd_ref:"CmdReqPlayerRefreshScore"`
	Score int32    `json:"score"`
}

// 玩家经营评分刷新响应
type ResPlayerRefreshScore struct { // 玩家经营评分刷新响应
	_    struct{} `cmd_ref:"CmdResPlayerRefreshScore"`
	Code int32    `json:"code"`
}

// 玩家编辑客户端数据请求
type ReqEditClientData struct {  
	_    struct{} `cmd_ref:"CmdReqEditClientData"`
	Data string   `json:"data"`
}

// 玩家编辑客户端数据响应
type ResEditClientData struct { 
	_    struct{} `cmd_ref:"CmdResEditClientData"`
	Code int32    `json:"code"`
}
