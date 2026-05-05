package route

import (
	"github.com/forfun/gforgame/examples/context"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	events "github.com/forfun/gforgame/examples/events"
	recharge "github.com/forfun/gforgame/examples/service/recharge"
	"github.com/forfun/gforgame/network"
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