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

	CmdHeroReqRecruit     = 5001
	CmdHeroResRecruit     = 5002
	CmdHeroResAllHero     = 5003
	CmdHeroReqLevelUp     = 5004
	CmdHeroResLevelUp     = 5005
	CmdHeroPushAdd        = 5006
	CmdHeroPushAttrChange = 5007
	CmdHeroReqUpStage     = 5008
	CmdHeroResUpStage     = 5009
)

type ReqPlayerLogin struct {
	Id  string
	Pwd string
}

type ReqPlayerLoadingFinish struct {
}

type ResPlayerLogin struct {
	Succ bool
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
	AttrType string  `json:"attrType"`
	Value    float32 `json:"value"`
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

type ReqPlayerUpStage struct {
}

type ResPlayerUpStage struct {
	Code int32 `json:"code"`
}
