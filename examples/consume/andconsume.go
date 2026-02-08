package consume

import (
	"fmt"
	"io/github/gforgame/examples/domain/player"
)

type AndConsume struct {
	Consumes []Consume
}

func (c *AndConsume) Add(consume Consume) {
	c.Consumes = append(c.Consumes, consume)
}

func (c *AndConsume) Verify(player *player.Player) error {
	for _, consume := range c.Consumes {
		if err := consume.Verify(player); err != nil {
			return err
		}
	}
	return nil
}

func (c *AndConsume) VerifySliently(player *player.Player) bool {
	err := c.Verify(player)
	return err == nil
}

func (c *AndConsume) Consume(player *player.Player, actionType int32) {
	for _, consume := range c.Consumes {
		consume.Consume(player, actionType)
	}
}


func (c *AndConsume) Merge() *AndConsume {
	merged := &AndConsume{}
	consumes := make(map[string]Consume)
	for _, consume := range c.Consumes {
		c.merge0(consumes, consume)
	}
	for _, consume := range consumes {
		merged.Add(consume)
	}
	return merged
}

func (c *AndConsume) merge0(consumes map[string]Consume, e Consume) {
	if andConsume, ok := e.(*AndConsume); ok {
		for _, consume := range andConsume.Consumes {
			c.merge0(consumes, consume)
		}
	} else if itemConsume, ok := e.(*ItemConsume); ok {
		key := fmt.Sprintf("item:%d", itemConsume.ItemId)
		if _, ok := consumes[key]; ok {
			consumes[key].(*ItemConsume).Amount += itemConsume.Amount
		} else {
			consumes[key] = itemConsume
		}
	} else if currencyConsume, ok := e.(*CurrencyConsume); ok {
		key := fmt.Sprintf("currency:%s", currencyConsume.Currency)
		if _, ok := consumes[key]; ok {
			consumes[key].(*CurrencyConsume).Amount += currencyConsume.Amount
		} else {
			consumes[key] = currencyConsume
		}
	}
}
