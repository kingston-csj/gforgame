package reward

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/player"
	"strconv"
)

type ItemReward struct {
	ItemId int32
	Amount int32
}

func (r *ItemReward) GetAmount() int {
	return int(r.Amount)
}

func (r *ItemReward) AddAmount(amount int) {
	r.Amount += int32(amount)
}
// 背包不限容量，所以不需要验证
func (r *ItemReward) Verify(player *player.Player) error {
	return nil
}

func (r *ItemReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *ItemReward) Reward(player *player.Player, actionType int) {
   itemOps := getItemOps()
	if itemOps == nil {
		return
	}
	if err := itemOps.AddByModelId(player, r.ItemId, r.Amount); err != nil {
		return
	}
}

func (r *ItemReward) GetType() string {
	return constants.RewardTypeItem
}

func (r *ItemReward) Serial() string {
	return strconv.Itoa(int(r.ItemId)) + "_" + strconv.Itoa(int(r.Amount))
}
