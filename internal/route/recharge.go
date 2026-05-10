package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	events "github.com/forfun/gforgame/internal/events"
	recharge "github.com/forfun/gforgame/internal/service/recharge"
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