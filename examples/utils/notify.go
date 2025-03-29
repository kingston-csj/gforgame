package utils

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/types"
)

func NotifyPlayer(player types.Player, event string, data interface{}) {
	s := context.SessionManager.GetSessionByPlayerId(player.GetID())
	if s == nil {
		return
	}
	s.SendWithoutIndex(data)
}
