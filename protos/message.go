package protos

const (
	CmdChatReqJoin = 1001
	CmdChatReqChat = 1002

	CmdPlayerReqLogin         = 2001
	CmdPlayerResLogin         = 2002
	CmdPlayerReqCreate        = 2003
	CmdPlayerResCreate        = 2004
	CmdPlayerReqLoadingFinish = 2005
	CmdPlayerReqUpLevel       = 2006
	CmdPlayerResUpLevel       = 2007
	CmdPlayerPushFightChange  = 2008
	CmdPlayerReqUpStage       = 2009
	CmdPlayerResUpStage       = 2010

	CmdGmReqAction = 3001
	CmdGmResAction = 3002

	CmdItemResBackpackInfo = 4001
	CmdItemResPurseInfo    = 4002
	CmdItemPushChanged     = 4003

	CmdHeroReqRecruit        = 5001
	CmdHeroResRecruit        = 5002
	CmdHeroResAllHero        = 5003
	CmdHeroReqLevelUp        = 5004
	CmdHeroResLevelUp        = 5005
	CmdHeroPushAdd           = 5006
	CmdHeroPushAttrChange    = 5007
	CmdHeroReqUpStage        = 5008
	CmdHeroResUpStage        = 5009
	CmdHeroReqCombine        = 5010
	CmdHeroResCombine        = 5011
	CmdHeroReqUpFight        = 5012
	CmdHeroResUpFight        = 5013
	CmdHeroReqOffFight       = 5014
	CmdHeroResOffFight       = 5015
	CmdHeroReqChangePosition = 5016
	CmdHeroResChangePosition = 5017

	CmdMailReqGetAllReward = 6001
	CmdMailResGetAllReward = 6002
	CmdMailReqRead         = 6003
	CmdMailResRead         = 6004
	CmdMailReqGetReward    = 6005
	CmdMailResGetReward    = 6006
	CmdMailReqDeleteAll    = 6007
	CmdMailResDeleteAll    = 6008
	CmdMailReqReadAll      = 6009
	CmdMailResReadAll      = 6010
	CmdMailPushAll         = 6011
)

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

type ReqHeroRecruit struct {
	Times int32 `json:"times"`
}

type ResHeroRecruit struct {
	Code        int32         `json:"code"`
	RewardInfos []*RewardInfo `json:"rewardInfos"`
}

type RewardInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PushPurseInfo struct {
	Diamond int32 `json:"diamond"`
	Gold    int32 `json:"gold"`
}

type ResAllHeroInfo struct {
	Heros []*HeroInfo `json:"heros"`
}

type HeroInfo struct {
	Id       int32      `json:"id"`
	Level    int32      `json:"level"`
	Position int32      `json:"position"`
	Stage    int32      `json:"stage"`
	Exp      int32      `json:"exp"`
	Hp       int32      `json:"hp"`
	Attrs    []AttrInfo `json:"attrs"`
	Fight    int32      `json:"fight"`
}

type ReqHeroLevelUp struct {
	HeroId  int32 `json:"heroId"`
	ToLevel int32 `json:"toLevel"`
}

type ResHeroLevelUp struct {
	Code int32 `json:"code"`
}

type ReqHeroUpStage struct {
	HeroId int32 `json:"heroId"`
}

type ResHeroUpStage struct {
	Code int32 `json:"code"`
}

type PushHeroAdd struct {
	HeroId int32 `json:"heroId"`
}

type PushItemChanged struct {
	ItemId int32 `json:"itemId"`
	Count  int32 `json:"count"`
}

type AttrInfo struct {
	AttrType string `json:"attrType"`
	Value    int32  `json:"value"`
}

type PushHeroAttrChange struct {
	HeroId int32      `json:"heroId"`
	Attrs  []AttrInfo `json:"attrs"`
	Fight  int32      `json:"fight"`
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

type ReqHeroCombine struct {
	HeroId int32 `json:"heroId"`
}

type ResHeroCombine struct {
	Code int32 `json:"code"`
}

type ReqHeroUpFight struct {
	HeroId   int32 `json:"heroId"`
	Position int32 `json:"position"`
}

type ResHeroUpFight struct {
	Code int32 `json:"code"`
}

type ReqHeroOffFight struct {
	HeroId int32 `json:"heroId"`
}

type ResHeroOffFight struct {
	Code int32 `json:"code"`
}

type ReqHeroChangePosition struct {
	HeroId   int32 `json:"heroId"`
	Position int32 `json:"position"`
}

type ResHeroChangePosition struct {
	Code  int32 `json:"code"`
	PosA  int32 `json:"posA"`
	HeroA int32 `json:"heroA"`
	PosB  int32 `json:"posB"`
	HeroB int32 `json:"heroB"`
}

type MailVo struct {
	Id int64 `json:"id"`
	// 邮件标题， 当TemplateId为0时，需要此字段
	Title string `json:"title"`
	// 邮件内容， 当TemplateId为0时，需要此字段
	Content string `json:"content"`
	// 邮件奖励
	Rewards []RewardInfo `json:"rewards"`
	// 邮件模板id
	TemplateId int32 `json:"templateId"`
	// 邮件状态
	Status int32 `json:"status"`
	// 邮件时间
	Time int64 `json:"time"`
}

type PushMailAll struct {
	Mails []MailVo `json:"mails"`
}

type ReqMailGetAllRewards struct{}

type ResMailGetAllRewards struct {
	Code int32 `json:"code"`
}

type ReqMailRead struct {
	Id int64 `json:"id"`
}

type ResMailRead struct {
	Code int32 `json:"code"`
}

type ReqMailGetReward struct {
	Id int64 `json:"id"`
}

type ResMailGetReward struct {
	Code int32 `json:"code"`
}

type ReqMailDeleteAll struct{}

type ResMailDeleteAll struct {
	Removed []int64 `json:"removed"`
}

type ResMailDelete struct {
	Code int32 `json:"code"`
}

type ReqMailReadAll struct{}

type ResMailReadAll struct {
	Code int32 `json:"code"`
}
