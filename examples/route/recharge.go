package route

import (
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
}