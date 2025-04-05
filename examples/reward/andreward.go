package reward

import "io/github/gforgame/examples/domain/player"

type AndReward struct {
	Rewards []Reward
}

func (r *AndReward) Verify(player *player.Player) error {
	for _, reward := range r.Rewards {
		if err := reward.Verify(player); err != nil {
			return err
		}
	}
	return nil
}

func (r *AndReward) VerifySliently(player *player.Player) bool {
	return r.Verify(player) == nil
}

func (r *AndReward) Reward(player *player.Player) {
	for _, reward := range r.Rewards {
		reward.Reward(player)
	}
}
