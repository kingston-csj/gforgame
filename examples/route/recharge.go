package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	events "io/github/gforgame/examples/events"
	recharge "io/github/gforgame/examples/service/recharge"
	"io/github/gforgame/network"
)

type RechargeRoute struct {
	network.Base
	service *recharge.RechargeService
}

func NewRechargeRoute() *RechargeRoute {
	return &RechargeRoute{}
}

func (c *RechargeRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		c.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}