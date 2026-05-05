package consume

import (
	"github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/examples/constants"
	"github.com/forfun/gforgame/examples/domain/player"
)

type ItemConsume struct {
	ItemId int32
	Amount int32
}

func (c *ItemConsume) Verify(player *player.Player) error {
	count := player.Backpack.GetItemCount(c.ItemId)
	if count < c.Amount {
		return errors.NewBusinessError(constants.I18N_ITEM_NOT_ENOUGH)
	}
	return nil
}

func (c *ItemConsume) VerifySliently(player *player.Player) bool {
	return c.Verify(player) == nil
}

func (c *ItemConsume) Consume(player *player.Player, actionType int32) {
	itemOps := GetItemOps()
	if itemOps == nil {
		return
	}
	if err := itemOps.UseByModelId(player.Id, c.ItemId, c.Amount); err != nil {
		return
	}
}
