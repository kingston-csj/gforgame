package io

import (
	"io/github/gforgame/examples/session"
	"io/github/gforgame/examples/types"
)

func NotifyPlayer(player types.Player, data any) {
	s := session.GetSessionByPlayerId(player.GetId())
	if s == nil {
		return
	}
	s.SendWithoutIndex(data)
}
