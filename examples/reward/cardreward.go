package reward

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/player"
	"strconv"
)

type CardReward struct {
	CardId int32
	Amount int32
}

func (r *CardReward) GetAmount() int {
	return int(r.Amount)
}

func (r *CardReward) AddAmount(amount int) {
	r.Amount += int32(amount)
}
// 背包不限容量，所以不需要验证
func (r *CardReward) Verify(player *player.Player) error {
	return nil
}

func (r *CardReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *CardReward) Reward(player *player.Player) {
    if ops := getItemOps(); ops != nil {
        // ops.AddByModelId(player, r.CardId, r.Amount)
    }
}

func (r *CardReward) GetType() string {
	return constants.RewardTypeCard
}

func (r *CardReward) Serial() string {
	return strconv.Itoa(int(r.CardId)) + "_" + strconv.Itoa(int(r.Amount))
}
