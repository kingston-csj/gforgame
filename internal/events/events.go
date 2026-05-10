package events

import "github.com/forfun/gforgame/internal/contract"

type IPlayerEvent interface {
	GetOwner() contract.Player
}

type PlayerEvent struct {
	Player contract.Player
}

func (e PlayerEvent) GetOwner() contract.Player {
	return e.Player
}

type LoginEvent struct {
	PlayerEvent
}

func (e LoginEvent) GetOwner() contract.Player {
	return e.Player
}

type RecruitEvent struct {
	PlayerEvent
	Times int32
	Type int32 // 招募类型 1:普通 2:高级
}

func (e RecruitEvent) GetOwner() contract.Player {
	return e.Player
}


type HeroGainEvent struct {
	PlayerEvent
	HeroId int32
}

func (e HeroGainEvent) GetOwner() contract.Player {
	return e.Player
}



type ItemConsumeEvent struct {
	PlayerEvent
	ItemId int32
	Count int32
}

func (e ItemConsumeEvent) GetOwner() contract.Player {
	return e.Player
}

type RechargeEvent struct {
	PlayerEvent
	RechargeId int32
}

func (e RechargeEvent) GetOwner() contract.Player {
	return e.Player
}

type ClientCustomEvent struct {
	PlayerEvent
	EventId int32
}

func (e ClientCustomEvent) GetOwner() contract.Player {
	return e.Player
}

type HeroEntrustEvent struct {
	PlayerEvent
}

func (e HeroEntrustEvent) GetOwner() contract.Player {
	return e.Player
}

type HeroLevelUpEvent struct {
	PlayerEvent
	HeroId int32
	Times int32
}
func (e HeroLevelUpEvent) GetOwner() contract.Player {
	return e.Player
}

type EquipLevelUpEvent struct {
	PlayerEvent
	EquipId int32
	Times int32
}
func (e EquipLevelUpEvent) GetOwner() contract.Player {
	return e.Player
}

type MallBuyEvent struct {
	PlayerEvent
}
func (e MallBuyEvent) GetOwner() contract.Player {
	return e.Player
}

type AreaScoreChangedEvent struct {
	PlayerEvent
	Score int32
}
func (e AreaScoreChangedEvent) GetOwner() contract.Player {
	return e.Player
}

type PassArenaEvent struct {
	PlayerEvent
}

func (e PassArenaEvent) GetOwner() contract.Player {
	return e.Player
}

type PassGuankaEvent struct {
	PlayerEvent
}

func (e PassGuankaEvent) GetOwner() contract.Player {
	return e.Player
}


type PassFubenEvent struct {
	PlayerEvent
}

func (e PassFubenEvent) GetOwner() contract.Player {
	return e.Player
}
