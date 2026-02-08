package events

import "io/github/gforgame/domain"

type IPlayerEvent interface {
	GetOwner() domain.Player
}

type PlayerEvent struct {
	Player domain.Player
}

func (e *PlayerEvent) GetOwner() domain.Player {
	return e.Player
}

type RecruitEvent struct {
	PlayerEvent
	Times int32
	Type int32 // 招募类型 1:普通 2:高级
}

type HeroGainEvent struct {
	PlayerEvent
	HeroId int32
}

type ItemConsumeEvent struct {
	PlayerEvent
	ItemId int32
	Count int32
}

type RechargeEvent struct {
	PlayerEvent
	RechargeId int32
}

type ClientCustomEvent struct {
	PlayerEvent
	EventId int32
}

type HeroEntrustEvent struct {
	PlayerEvent
}

type HeroLevelUpEvent struct {
	PlayerEvent
	HeroId int32
	Times int32
}

type MallBuyEvent struct {
	PlayerEvent
}

type AreaScoreChangedEvent struct {
	PlayerEvent
	Score int32
}

type PassArenaEvent struct {
	PlayerEvent
}
