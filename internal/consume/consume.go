package consume

import "github.com/forfun/gforgame/internal/domain/player"

type Consume interface {
	Verify(player *player.Player) error

	VerifySliently(player *player.Player) bool

	Consume(player *player.Player, actionType int32)
}
