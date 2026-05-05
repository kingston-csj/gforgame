package reward

import (
	"strconv"

	"github.com/forfun/gforgame/examples/constants"
	"github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/io"
	"github.com/forfun/gforgame/examples/protos"
)

type CurrencyReward struct {
	Currency   string
	Amount int32
}

func NewCurrencyReward(kind string, amount int32) *CurrencyReward {
	return &CurrencyReward{
		Currency:   kind,
		Amount: amount,
	}
}

func (r *CurrencyReward) GetAmount() int {
	return int(r.Amount)
}

func (r *CurrencyReward) AddAmount(amount int) {
	r.Amount += int32(amount)
}

func (r *CurrencyReward) Verify(player *player.Player) error {
	return nil
}

func (r *CurrencyReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *CurrencyReward) Reward(player *player.Player, actionType int) {
    if ops := getCurrencyOps(); ops != nil {
        ops.Add(player.Id, r.Currency, r.Amount)
        return
    }
    if r.Currency == "gold" {
        player.Purse.AddGold(r.Amount)
    } else if r.Currency == "diamond" {
        player.Purse.AddDiamond(r.Amount)
    }
    io.NotifyPlayer(player, &protos.PushPurseInfo{
        Gold:    player.Purse.Gold,
        Diamond: player.Purse.Diamond,
    })
}

func (r *CurrencyReward) GetType() string {
	return constants.RewardTypeCurrency
}

func (r *CurrencyReward) Serial() string {
	return r.Currency + "_" + strconv.Itoa(int(r.Amount))
}


