package route

import (
	recharge "github.com/forfun/gforgame/internal/service/recharge"
)

type RechargeRoute struct {
	service *recharge.RechargeService
}

func NewRechargeRoute() *RechargeRoute {
	return &RechargeRoute{}
}
