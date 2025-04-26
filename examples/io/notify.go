package io

import (
	"io/github/gforgame/domain"
	"io/github/gforgame/network"
)

func NotifyPlayer(player domain.Player, data any) {
	s := network.GetSessionByPlayerId(player.GetId())
	if s == nil {
		return
	}
	s.SendWithoutIndex(data)
}
