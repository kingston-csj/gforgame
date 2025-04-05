package consume

import "io/github/gforgame/examples/domain/player"

type AndConsume struct {
	Consumes []Consume
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

func (c *AndConsume) Consume(player *player.Player) {
	for _, consume := range c.Consumes {
		consume.Consume(player)
	}
}
