package reward

import (
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/item"
)

type ItemReward struct {
	ItemId int32
	Amount int32
}

// 背包不限容量，所以不需要验证
func (r *ItemReward) Verify(player *player.Player) error {
	return nil
}

func (r *ItemReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *ItemReward) Reward(player *player.Player) {
	item.GetItemService().AddByModelId(player, r.ItemId, r.Amount)
}
