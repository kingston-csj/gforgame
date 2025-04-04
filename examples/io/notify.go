package io

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/types"
)

func NotifyPlayer(player types.Player, data any) {
	s := context.SessionManager.GetSessionByPlayerId(player.GetID())
	if s == nil {
		return
	}
	s.SendWithoutIndex(data)
}
