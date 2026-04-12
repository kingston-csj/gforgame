package io

import (
	"io/github/gforgame/examples/contract"
	"io/github/gforgame/network"
)

func NotifyPlayer(player contract.Player, data any) {
	s := network.GetSessionByPlayerId(player.GetId())
	if s == nil {
		return
	}
	s.SendWithoutIndex(data)
}
