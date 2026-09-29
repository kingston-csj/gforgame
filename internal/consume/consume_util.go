package consume

import (
	"strings"

	"github.com/forfun/gforgame/common/util/conv"
)

func ParseConsume(config string) *AndConsume {
	if conv.IsBlankString(config) {
		return &AndConsume{}
	}
	splits := strings.Split(config, ",")
	andConsume := &AndConsume{}
	for _, split := range splits {
		params := strings.Split(split, "_")
		consumeType := params[0]
		if conv.EqualsIgnoreCase(consumeType, "item") {
			itemId, _ := conv.StringToInt32(params[1])
			count, _ := conv.StringToInt32(params[2])
			itemConsume := &ItemConsume{
				ItemId: itemId,
				Amount: count,
			}
			andConsume.Add(itemConsume)
		} else if conv.EqualsIgnoreCase(consumeType, "currency") {
			amount, _ := conv.StringToInt32(params[2])
			currencyConsume := &CurrencyConsume{
				Currency: params[1],
				Amount:   amount,
			}
			andConsume.Add(currencyConsume)
		}
	}
	return andConsume
}
