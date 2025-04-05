package reward

import (
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

type CurrencyReward struct {
	Kind   string
	Amount int32
}

func (r *CurrencyReward) Verify(player *player.Player) error {
	return nil
}

func (r *CurrencyReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *CurrencyReward) Reward(player *player.Player) {
	if r.Kind == "gold" {
		player.Purse.AddGold(r.Amount)
	} else if r.Kind == "diamond" {
		player.Purse.AddDiamond(r.Amount)
	}
	io.NotifyPlayer(player, &protos.PushPurseInfo{
		Gold:    player.Purse.Gold,
		Diamond: player.Purse.Diamond,
	})
}
