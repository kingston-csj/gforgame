package item

import "io/github/gforgame/examples/domain/player"

type ItemConsumeOps interface {
	UseByModelId(player *player.Player, itemId int32, amount int32) error
}
