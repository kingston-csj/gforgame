package consume

import (
	"strings"

	"github.com/forfun/gforgame/common/util/conv"
)

func ParseConsume(config string) *AndConsume{
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
				Amount:  count,
			}
			andConsume.Add(itemConsume)
		} else if conv.EqualsIgnoreCase(consumeType, "Gold") {
			amount, _ := conv.StringToInt32(params[1])
			currencyConsume := &CurrencyConsume{
				Currency: "Gold",	
				Amount:   amount,
			}
			andConsume.Add(currencyConsume)
		} else if conv.EqualsIgnoreCase(consumeType, "Diamond") {
			amount, _ := conv.StringToInt32(params[1])
			currencyConsume := &CurrencyConsume{
				Currency: "Diamond",	
				Amount:   amount,
			}
			andConsume.Add(currencyConsume)
		}
	}
	return andConsume
}