package reward

import (
	"strconv"

	"github.com/forfun/gforgame/examples/constants"
	"github.com/forfun/gforgame/examples/domain/player"
)

type RuneReward struct {
	ItemId int32
	Amount int32
}

func (r *RuneReward) GetAmount() int {
	return int(r.Amount)
}

func (r *RuneReward) AddAmount(amount int) {
	r.Amount += int32(amount)
}
// 背包不限容量，所以不需要验证
func (r *RuneReward) Verify(player *player.Player) error {
	return nil
}

func (r *RuneReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *RuneReward) Reward(player *player.Player, actionType int) {
    if ops := getItemOps(); ops != nil {
        ops.AddByModelId(player.Id, r.ItemId, r.Amount)
    }
}

func (r *RuneReward) GetType() string {
	return constants.RewardTypeRune
}

func (r *RuneReward) Serial() string {
	return strconv.Itoa(int(r.ItemId)) + "_" + strconv.Itoa(int(r.Amount))
}
