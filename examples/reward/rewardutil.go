package reward

import (
	"errors"
	"io/github/gforgame/domain"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
	"math"
	"strings"
)

func ParseReward(config string) *AndReward {
	if util.IsBlankString(config) {
		return NewAndReward()
	}
	rewards := ParseRewardList(config)
	andReward := NewAndReward()
	for _, reward := range rewards {
		andReward.AddReward(reward)
	}
	return andReward
}

func ParseRewardList(config string) []Reward {
	result := make([]Reward, 0)
	if util.IsBlankString(config) {
		return result
	}
	splits := strings.Split(config, ",")
	for _, split := range splits {
		params := strings.Split(split, "_")
		rewardType := params[0]
		if util.EqualsIgnoreCase(rewardType, constants.CurrencyTypeGold) {
			amount, _ := util.StringToInt32(params[1])
			result = append(result, &CurrencyReward{
				Currency:   "Gold",
				Amount:     amount,
			})
		} else if util.EqualsIgnoreCase(rewardType, constants.CurrencyTypeDiamond) {
			amount, _ := util.StringToInt32(params[1])
			result = append(result, &CurrencyReward{
				Currency:   "Diamond", 
				Amount:     amount,
			})
		} else if util.EqualsIgnoreCase(rewardType, constants.RewardTypeItem) {
			itemId, _ := util.StringToInt32(params[1])
			amount, _ := util.StringToInt32(params[2])
			result = append(result, &ItemReward{
				ItemId: itemId,
				Amount: amount,
			})
		} else if util.EqualsIgnoreCase(rewardType, constants.RewardTypeRune) {
			amount, _ := util.StringToInt32(params[1])
			result = append(result, &RuneReward{
				ItemId: amount,
			})
		} else if util.EqualsIgnoreCase(rewardType, constants.RewardTypeTicket) {
			mapId, _ := util.StringToInt32(params[1])
			amount, _ := util.StringToInt32(params[2])
			result = append(result, &TicketReward{
				MapId: mapId,
				Amount: amount,
			})
		} else if util.EqualsIgnoreCase(rewardType, constants.RewardTypeHero) {
			amount, _ := util.StringToInt32(params[1])
			result = append(result, &HeroReward{
				HeroId: amount,
			})
		} else if util.EqualsIgnoreCase(rewardType, constants.RewardTypeCard) {
			cardId, _ := util.StringToInt32(params[1])
			amount, _ := util.StringToInt32(params[2])
			result = append(result, &CardReward{
				CardId: cardId,
				Amount: amount,
			})
		}
	}
	return result
}


func ParseRewards(rewards []domain.RewardDefLite) *AndReward {
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
            andReward.AddReward(&CurrencyReward{Currency: kind, Amount: amount})
        }
    }
    return andReward
}

func FromConsumes(consumes []consume.Consume) *AndReward {
	andReward := NewAndReward()
	for _, c := range consumes {
		reward, err := consume2Reward(c)
		if err != nil {
			continue
		}
		andReward.AddReward(reward)
	}
	return andReward
}

func consume2Reward(c consume.Consume) (Reward, error) {
    // 如果类型是 CurrencyConsume
    if currencyConsume, ok := c.(*consume.CurrencyConsume); ok {
        return &CurrencyReward{
            Currency: currencyConsume.Currency,
            Amount:   currencyConsume.Amount,
        }, nil
    } else if itemConsume, ok := c.(*consume.ItemConsume); ok {
        return &ItemReward{
            ItemId: itemConsume.ItemId,
            Amount: itemConsume.Amount,
        }, nil
    }
    return nil, errors.New("unsupported consume type")
}

func ToRewardVos(rewards Reward) []*protos.RewardVo {
	rewardVos := make([]*protos.RewardVo, 0)
	switch rewards.(type) {
	case *AndReward:
		andReward := rewards.(*AndReward)
		for _, reward := range andReward.Rewards {
			rewardVos = append(rewardVos, ToRewardVo(reward))
		}
	default:
		rewardVos = append(rewardVos, ToRewardVo(rewards))
	}
	return rewardVos
}

func ToRewardVo(reward Reward) *protos.RewardVo {
	return &protos.RewardVo{
		Type:  reward.GetType(),
		Value: reward.Serial(),
	}
}

func MultiplyAndReward(sourceRewards *AndReward, multiple float64) *AndReward {
	andReward := NewAndReward()
	for _, reward := range sourceRewards.Rewards {
		// 数量全部向上取整
		rewardAmount := int32(math.Ceil(float64(reward.GetAmount()) * multiple))
		andReward.AddReward(modifyRewardAmount(reward, rewardAmount));
	}
	return andReward;
}

func modifyRewardAmount(reward Reward, amount int32) Reward {
	if _, ok := reward.(*CurrencyReward); ok {
		return &CurrencyReward{Currency: reward.(*CurrencyReward).Currency, Amount: amount}
	} else if _, ok := reward.(*ItemReward); ok {
		return &ItemReward{ItemId: reward.(*ItemReward).ItemId, Amount: amount}
	} else if _, ok := reward.(*TicketReward); ok {
		return &TicketReward{MapId: reward.(*TicketReward).MapId, Amount: amount}
	} else if _, ok := reward.(*RuneReward); ok {
		return &RuneReward{ItemId: reward.(*RuneReward).ItemId, Amount: amount}
	} else if _, ok := reward.(*HeroReward); ok {
		return &HeroReward{HeroId: reward.(*HeroReward).HeroId, Amount: amount}
	} else if _, ok := reward.(*CardReward); ok {
		return &CardReward{CardId: reward.(*CardReward).CardId, Amount: amount}
	} else {
		panic("unsupported reward type: " + reward.GetType())
	}
}

// 奖励加倍（数量为向上取整）
func multiply(sourceRewards Reward, multiple float64) Reward {
	if _, ok := sourceRewards.(*AndReward); ok {
		return MultiplyAndReward(sourceRewards.(*AndReward), multiple);
	} else {
		return modifyRewardAmount(sourceRewards, int32(math.Ceil(float64(sourceRewards.GetAmount()) * multiple)));
	}
}
