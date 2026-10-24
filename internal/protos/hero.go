package protos

type ReqHeroRecruit struct { // 英雄——招募
	_       struct{} `cmd_ref:"CmdHeroReqRecruit" type:"req"`
	Counter int32    `json:"counter"`
}

type ResHeroRecruit struct { // 英雄——招募
	_         struct{}    `cmd_ref:"CmdHeroResRecruit" type:"res"`
	Code      int32       `json:"code"`
	RewardVos []*RewardVo `json:"rewardVos"`
}

type PushAllHeroInfo struct { // 英雄——推送所有英雄信息
	_     struct{}    `cmd_ref:"CmdHeroPushAllHero" type:"push"`
	Heros []*HeroInfo `json:"heros"`
}

type AttrInfo struct { // 英雄属性信息
	AttrType string `json:"attrType"`
	Value    int32  `json:"value"`
}

type HeroInfo struct { // 英雄信息
	Id       int32      `json:"id"`
	Level    int32      `json:"level"`
	Position int32      `json:"position"`
	Stage    int32      `json:"stage"`
	Exp      int32      `json:"exp"`
	Attrs    []AttrInfo `json:"attrs"`
	Fight    int32      `json:"fight"`
	Equips   []ItemInfo `json:"equips"`
}

type ReqHeroLevelUp struct { // 英雄——升级
	_       struct{} `cmd_ref:"CmdHeroReqLevelUp" type:"req"`
	HeroId  int32    `json:"heroId"`
	ToLevel int32    `json:"toLevel"`
}

type ResHeroLevelUp struct { // 英雄——升级
	_    struct{} `cmd_ref:"CmdHeroResLevelUp" type:"res"`
	Code int32    `json:"code"`
}

type ReqHeroUpStage struct { // 英雄——升级阶段
	_      struct{} `cmd_ref:"CmdHeroReqUpStage" type:"req"`
	HeroId int32    `json:"heroId"`
}

type ResHeroUpStage struct { // 英雄——升级阶段
	_    struct{} `cmd_ref:"CmdHeroResUpStage" type:"res"`
	Code int32    `json:"code"`
}

// 英雄——添加
type PushHeroAdd struct {
	_      struct{} `cmd_ref:"CmdHeroPushAdd" type:"push"`
	HeroId int32    `json:"heroId"`
}

// 英雄——属性变化
type PushHeroAttrChange struct {
	_      struct{}   `cmd_ref:"CmdHeroPushAttrChange" type:"push"`
	HeroId int32      `json:"heroId"`
	Attrs  []AttrInfo `json:"attrs"`
	Fight  int32      `json:"fight"`
}

// 英雄——合成
type ReqHeroCombine struct {
	_      struct{} `cmd_ref:"CmdHeroReqCombine" type:"req"`
	HeroId int32    `json:"heroId"`
}

// 英雄——合成
type ResHeroCombine struct {
	_    struct{} `cmd_ref:"CmdHeroResCombine" type:"res"`
	Code int32    `json:"code"`
}

// 英雄——上阵
type ReqHeroUpFight struct {
	_        struct{} `cmd_ref:"CmdHeroReqUpFight" type:"req"`
	HeroId   int32    `json:"heroId"`
	Position int32    `json:"position"`
}

// 英雄——上阵
type ResHeroUpFight struct {
	_    struct{} `cmd_ref:"CmdHeroResUpFight" type:"res"`
	Code int32    `json:"code"`
}

// 英雄——下阵
type ReqHeroOffFight struct {
	_      struct{} `cmd_ref:"CmdHeroReqOffFight" type:"req"`
	HeroId int32    `json:"heroId"`
}

// 英雄——下阵
type ResHeroOffFight struct {
	_    struct{} `cmd_ref:"CmdHeroResOffFight" type:"res"`
	Code int32    `json:"code"`
}

// 英雄——交换位置
type ReqHeroChangePosition struct {
	_        struct{} `cmd_ref:"CmdHeroReqChangePosition" type:"req"`
	HeroId   int32    `json:"heroId"`
	Position int32    `json:"position"`
}

// 英雄——交换位置
type ResHeroChangePosition struct {
	_     struct{} `cmd_ref:"CmdHeroResChangePosition" type:"res"`
	Code  int32    `json:"code"`
	PosA  int32    `json:"posA"`
	HeroA int32    `json:"heroA"`
	PosB  int32    `json:"posB"`
	HeroB int32    `json:"heroB"`
}
