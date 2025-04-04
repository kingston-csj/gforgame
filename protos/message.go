package protos

const (
	CmdChatReqJoin = 1001
	CmdChatReqChat = 1002

	CmdPlayerReqLogin  = 2001
	CmdPlayerResLogin  = 2002
	CmdPlayerReqCreate = 2003
	CmdPlayerResCreate = 2004

	CmdGmReqAction = 3001
	CmdGmResAction = 3002

	CmdItemResBackpackInfo = 4001
	CmdItemResPurseInfo    = 4002

	CmdHeroReqRecruit = 5001
	CmdHeroResRecruit = 5002
	CmdHeroResAllHero = 5003
	CmdHeroReqLevelUp = 5004
	CmdHeroResLevelUp = 5005
	CmdHeroPushAdd    = 5006
)

type ReqPlayerLogin struct {
	Id  string
	Pwd string
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
	Id       int32 `json:"id"`
	Level    int32 `json:"level"`
	Position int32 `json:"position"`
	Stage    int32 `json:"stage"`
}

type ReqHeroLevelUp struct {
	HeroId  int32 `json:"heroId"`
	ToLevel int32 `json:"toLevel"`
}

type ResHeroLevelUp struct {
	Code int32 `json:"code"`
}

type PushHeroAdd struct {
	HeroId int32 `json:"heroId"`
}
