package item

import (
	playerdomain "io/github/gforgame/examples/domain/player"
)

type ItemRewardOps interface {
    AddByModelId(p *playerdomain.Player, itemId int32, amount int32) error
}

type CurrencyOps interface {
    Add(p *playerdomain.Player, kind string, amount int32)
}
