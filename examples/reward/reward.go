package reward

import "io/github/gforgame/examples/domain/player"

type Reward interface {
	Verify(player *player.Player) error

	VerifySliently(player *player.Player) bool

	Reward(player *player.Player)
}
