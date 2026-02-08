package consume

import (
	"io/github/gforgame/util"
	"strings"
)

func ParseConsume(config string) *AndConsume{
	if util.IsBlankString(config) {
		return &AndConsume{}
	}
	splits := strings.Split(config, ",")
	andConsume := &AndConsume{}
	for _, split := range splits {
		params := strings.Split(split, "_")
		consumeType := params[0]
		if util.EqualsIgnoreCase(consumeType, "item") {
			itemId, _ := util.StringToInt32(params[1])
			count, _ := util.StringToInt32(params[2])
			itemConsume := &ItemConsume{
				ItemId: itemId,
				Amount:  count,
			}
			andConsume.Add(itemConsume)
		} else if util.EqualsIgnoreCase(consumeType, "Gold") {
			amount, _ := util.StringToInt32(params[1])
			currencyConsume := &CurrencyConsume{
				Currency: "Gold",	
				Amount:   amount,
			}
			andConsume.Add(currencyConsume)
		} else if util.EqualsIgnoreCase(consumeType, "Diamond") {
			amount, _ := util.StringToInt32(params[1])
			currencyConsume := &CurrencyConsume{
				Currency: "Diamond",	
				Amount:   amount,
			}
			andConsume.Add(currencyConsume)
		}
	}
	return andConsume
}