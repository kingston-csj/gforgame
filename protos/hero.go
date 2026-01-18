package protos

type ReqHeroRecruit struct {
	_       struct{} `cmd_ref:"CmdHeroReqRecruit" type:"req"`
	Counter int32    `json:"counter"`
}

type ResHeroRecruit struct {
	_         struct{}    `cmd_ref:"CmdHeroResRecruit" type:"res"`
	Code      int32       `json:"code"`
	RewardVos []*RewardVo `json:"rewardVos"`
}

type PushAllHeroInfo struct {
	_     struct{}    `cmd_ref:"CmdHeroPushAllHero" type:"push"`
	Heros []*HeroInfo `json:"heros"`
}

type AttrInfo struct {
	AttrType string `json:"attrType"`
	Value    int32  `json:"value"`
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
	_       struct{} `cmd_ref:"CmdHeroReqLevelUp" type:"req"`
	HeroId  int32    `json:"heroId"`
	ToLevel int32    `json:"toLevel"`
}

type ResHeroLevelUp struct {
	_    struct{} `cmd_ref:"CmdHeroResLevelUp" type:"res"`
	Code int32    `json:"code"`
}

type ReqHeroUpStage struct {
	_      struct{} `cmd_ref:"CmdHeroReqUpStage" type:"req"`
	HeroId int32    `json:"heroId"`
}

type ResHeroUpStage struct {
	_    struct{} `cmd_ref:"CmdHeroResUpStage" type:"res"`
	Code int32    `json:"code"`
}

type PushHeroAdd struct {
	_      struct{} `cmd_ref:"CmdHeroPushAdd" type:"push"`
	HeroId int32    `json:"heroId"`
}

type PushHeroAttrChange struct {
	_      struct{}   `cmd_ref:"CmdHeroPushAttrChange" type:"push"`
	HeroId int32      `json:"heroId"`
	Attrs  []AttrInfo `json:"attrs"`
	Fight  int32      `json:"fight"`
}

type ReqHeroCombine struct {
	_      struct{} `cmd_ref:"CmdHeroReqCombine" type:"req"`
	HeroId int32    `json:"heroId"`
}

type ResHeroCombine struct {
	_    struct{} `cmd_ref:"CmdHeroResCombine" type:"res"`
	Code int32    `json:"code"`
}

type ReqHeroUpFight struct {
	_        struct{} `cmd_ref:"CmdHeroReqUpFight" type:"req"`
	HeroId   int32    `json:"heroId"`
	Position int32    `json:"position"`
}

type ResHeroUpFight struct {
	_    struct{} `cmd_ref:"CmdHeroResUpFight" type:"res"`
	Code int32    `json:"code"`
}

type ReqHeroOffFight struct {
	_      struct{} `cmd_ref:"CmdHeroReqOffFight" type:"req"`
	HeroId int32    `json:"heroId"`
}

type ResHeroOffFight struct {
	_    struct{} `cmd_ref:"CmdHeroResOffFight" type:"res"`
	Code int32    `json:"code"`
}

type ReqHeroChangePosition struct {
	_        struct{} `cmd_ref:"CmdHeroReqChangePosition" type:"req"`
	HeroId   int32    `json:"heroId"`
	Position int32    `json:"position"`
}

type ResHeroChangePosition struct {
	_     struct{} `cmd_ref:"CmdHeroResChangePosition" type:"res"`
	Code  int32    `json:"code"`
	PosA  int32    `json:"posA"`
	HeroA int32    `json:"heroA"`
	PosB  int32    `json:"posB"`
	HeroB int32    `json:"heroB"`
}
