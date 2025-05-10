package protos

type ReqHeroRecruit struct {
	Times int32 `json:"times"`
}

type ResHeroRecruit struct {
	Code        int32         `json:"code"`
	RewardInfos []*RewardInfo `json:"rewardInfos"`
}

type ResAllHeroInfo struct {
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

type PushHeroAttrChange struct {
	HeroId int32      `json:"heroId"`
	Attrs  []AttrInfo `json:"attrs"`
	Fight  int32      `json:"fight"`
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
