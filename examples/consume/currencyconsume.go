package consume

import (
	"io/github/gforgame/common"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

type CurrencyConsume struct {
	Kind   string
	Amount int32
}

func (c *CurrencyConsume) Verify(player *player.Player) error {
	purse := player.Purse
	if c.Kind == "gold" {
		if !purse.IsEnoughGold(c.Amount) {
			return common.NewBusinessRequestException(constants.Gold_NOT_ENOUGH)
		}
	} else if c.Kind == "diamond" {
		if !purse.IsEnoughDiamond(c.Amount) {
			return common.NewBusinessRequestException(constants.Diamond_NOT_ENOUGH)
		}
	}

	return nil
}

func (c *CurrencyConsume) VerifySliently(player *player.Player) bool {
	return c.Verify(player) == nil
}

func (c *CurrencyConsume) Consume(player *player.Player) {
	if c.Kind == "gold" {
		player.Purse.SubGold(c.Amount)
	} else if c.Kind == "diamond" {
		player.Purse.SubDiamond(c.Amount)
	}
	io.NotifyPlayer(player, &protos.PushPurseInfo{
		Gold:    player.Purse.Gold,
		Diamond: player.Purse.Diamond,
	})
}
