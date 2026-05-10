package io

import (
	"github.com/forfun/gforgame/internal/contract"
	"github.com/forfun/gforgame/network"
)

func NotifyPlayer(player contract.Player, data any) {
	s := network.GetSessionByPlayerId(player.GetId())
	if s == nil {
		return
	}
	s.SendWithoutIndex(data)
}
