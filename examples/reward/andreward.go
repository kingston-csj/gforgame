package reward

import (
	"fmt"

	"io/github/gforgame/examples/domain/player"
)

type AndReward struct {
	Rewards []Reward
}

func NewAndReward() *AndReward {
	return &AndReward{
		Rewards: make([]Reward, 0),
	}
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

func (r *AndReward) AddReward(reward Reward) {
	r.Rewards = append(r.Rewards, reward)
}

func (a *AndReward) Merge() *AndReward {
	result := make(map[string]Reward)
	for _, r := range a.Rewards {
		merge0(result, r)
	}
	merged := &AndReward{}
	for _, v := range result {
		merged.Rewards = append(merged.Rewards, v)
	}
	return merged
}

func merge0(result map[string]Reward, r Reward) {
	switch v := r.(type) {
	case *AndReward:
		for _, child := range v.Rewards {
			merge0(result, child)
		}
	case *CurrencyReward:
		key := "currency:" + v.Kind
		if prev, ok := result[key]; ok {
			prevMoney := prev.(*CurrencyReward)
			result[key] = &CurrencyReward{
				Kind:   v.Kind,
				Amount: prevMoney.Amount + v.Amount,
			}
		} else {
			result[key] = &CurrencyReward{
				Kind:   v.Kind,
				Amount: v.Amount,
			}
		}
	case *ItemReward:
		key := "item:" + fmt.Sprintf("%d", v.ItemId)
		if prev, ok := result[key]; ok {
			prevItem := prev.(*ItemReward)
			result[key] = &ItemReward{
				ItemId: v.ItemId,
				Amount: prevItem.Amount + v.Amount,
			}
		} else {
			result[key] = &ItemReward{
				ItemId: v.ItemId,
				Amount: v.Amount,
			}
		}

	default:
		panic(fmt.Sprintf("cannot merge %T", r))
	}
}
