package events

import "io/github/gforgame/domain"

type RecruitEvent struct {
	Player domain.Player
	Times int32
}

type HeroGainEvent struct {
	Player domain.Player
	HeroId int32
}

type ItemConsumeEvent struct {
	Player domain.Player
	ItemId int32
	Count int32
}
