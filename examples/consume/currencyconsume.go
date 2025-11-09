package consume

import (
	"io/github/gforgame/common"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

type CurrencyConsume struct {
	Currency string
	Amount   int32
}

func (c *CurrencyConsume) Verify(player *player.Player) error {
	purse := player.Purse
	if c.Currency == "gold" {
		if !purse.IsEnoughGold(c.Amount) {
			return common.NewBusinessRequestException(constants.I18N_GOLD_NOT_ENOUGH)
		}
	} else if c.Currency == "diamond" {	
		if !purse.IsEnoughDiamond(c.Amount) {
			return common.NewBusinessRequestException(constants.I18N_DIAMOND_NOT_ENOUGH)
		}
	}

	return nil
}

func (c *CurrencyConsume) VerifySliently(player *player.Player) bool {
	return c.Verify(player) == nil
}

func (c *CurrencyConsume) Consume(player *player.Player) {
	if c.Currency == "gold" {
		player.Purse.SubGold(c.Amount)
	} else if c.Currency == "diamond" {
		player.Purse.SubDiamond(c.Amount)
	}
	io.NotifyPlayer(player, &protos.PushPurseInfo{
		Gold:    player.Purse.Gold,
		Diamond: player.Purse.Diamond,
	})
}
