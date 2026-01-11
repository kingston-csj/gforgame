package reward

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/player"
	"strconv"
)

type TicketReward struct {
	MapId int32
	Amount int32
}

func (r *TicketReward) GetAmount() int {
	return int(r.Amount)
}

func (r *TicketReward) AddAmount(amount int) {
	r.Amount += int32(amount)
}
// 背包不限容量，所以不需要验证
func (r *TicketReward) Verify(player *player.Player) error {
	return nil
}

func (r *TicketReward) VerifySliently(player *player.Player) bool {
	return true
}

func (r *TicketReward) Reward(player *player.Player, actionType int) {
    if ops := getItemOps(); ops != nil {
        ops.AddByModelId(player, r.MapId, r.Amount)
    }
}

func (r *TicketReward) GetType() string {
	return constants.RewardTypeTicket
}

func (r *TicketReward) Serial() string {
	return strconv.Itoa(int(r.MapId)) + "_" + strconv.Itoa(int(r.Amount))
}
