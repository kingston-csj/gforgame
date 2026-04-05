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
	Type  int32 // 招募类型 1:普通 2:高级
}

func (e *RecruitEvent) GetOwner() domain.Player {
	return e.Player
}

type HeroGainEvent struct {
	PlayerEvent
	HeroId int32
}

func (e *HeroGainEvent) GetOwner() domain.Player {
	return e.Player
}

type ItemConsumeEvent struct {
	PlayerEvent
	ItemId int32
	Count  int32
}

func (e *ItemConsumeEvent) GetOwner() domain.Player {
	return e.Player
}

type RechargeEvent struct {
	PlayerEvent
	RechargeId int32
}

func (e *RechargeEvent) GetOwner() domain.Player {
	return e.Player
}

type ClientCustomEvent struct {
	PlayerEvent
	EventId int32
}

func (e *ClientCustomEvent) GetOwner() domain.Player {
	return e.Player
}

type HeroEntrustEvent struct {
	PlayerEvent
}

func (e *HeroEntrustEvent) GetOwner() domain.Player {
	return e.Player
}

type HeroLevelUpEvent struct {
	PlayerEvent
	HeroId int32
	Times  int32
}

func (e *HeroLevelUpEvent) GetOwner() domain.Player {
	return e.Player
}

type MallBuyEvent struct {
	PlayerEvent
}

func (e *MallBuyEvent) GetOwner() domain.Player {
	return e.Player
}

type AreaScoreChangedEvent struct {
	PlayerEvent
	Score int32
}

func (e *AreaScoreChangedEvent) GetOwner() domain.Player {
	return e.Player
}

type PassArenaEvent struct {
	PlayerEvent
}

func (e *PassArenaEvent) GetOwner() domain.Player {
	return e.Player
}
