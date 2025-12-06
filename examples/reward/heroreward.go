package reward

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/player"
	"strconv"
)

type HeroReward struct {
	HeroId int32
	Amount int32
}

func (r *HeroReward) GetAmount() int {
	return int(r.Amount)
}

func (r *HeroReward) AddAmount(amount int) {
	r.Amount += int32(amount)
}
// 背包不限容量，所以不需要验证
func (r *HeroReward) Verify(player *player.Player) error {
	return nil
}

func (r *HeroReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *HeroReward) Reward(player *player.Player) {
    if ops := getItemOps(); ops != nil {
        ops.AddByModelId(player, r.HeroId, r.Amount)
    }
}

func (r *HeroReward) GetType() string {
    return constants.RewardTypeHero
}

func (r *HeroReward) Serial() string {
	return strconv.Itoa(int(r.HeroId)) + "_" + strconv.Itoa(int(r.Amount))
}
