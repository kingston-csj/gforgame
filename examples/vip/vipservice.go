package vip

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	"sync"
) 


type VipService struct{}


var (
	instance *VipService
	once     sync.Once
)

func GetVipService() *VipService {
	once.Do(func() {
		instance = &VipService{}	
	})
	return instance
}

func (v *VipService) CheckRecharge(p *playerdomain.Player, rechargeId int32 ) {
	
}