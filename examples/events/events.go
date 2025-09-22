package events

import "io/github/gforgame/domain"

type RecruitEvent struct {
	Player domain.Player
	Times int32
}