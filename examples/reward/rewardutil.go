package reward

import (
	"strings"

	"io/github/gforgame/examples/domain/config/item"
	"io/github/gforgame/util"
)

func ParseRewards(rewards []item.RewardDef) *AndReward {
	andReward := NewAndReward()

	for _, rewardItem := range rewards {
		split := strings.Split(rewardItem.Value, "=")
		switch rewardItem.Type {
		case "item":
			itemId, _ := util.StringToInt32(split[0])
			amount, _ := util.StringToInt32(split[1])
			andReward.AddReward(&ItemReward{
				ItemId: itemId,
				Amount: amount,
			})
		case "currency":
			kind := split[0]
			amount, _ := util.StringToInt32(split[1])
			andReward.AddReward(&CurrencyReward{
				Kind:   kind,
				Amount: amount,
			})
		}
	}
	return andReward
}
