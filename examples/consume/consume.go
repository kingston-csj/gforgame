package consume

import "io/github/gforgame/examples/domain/player"

type Consume interface {
	Verify(player *player.Player) error

	VerifySliently(player *player.Player) bool

	Consume(player *player.Player, actionType int32)
}
