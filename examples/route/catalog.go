package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/catalog"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type CatalogRoute struct {
	network.Base
	service *catalog.CatalogService
}

func NewCatalogRoute() *CatalogRoute {
	return &CatalogRoute{
		service: catalog.GetCatalogService(),
	}
}

func (ps *CatalogRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *CatalogRoute) ReqCatalogReward(s *network.Session, index int32, msg *protos.ReqCatalogReward) *protos.ResCatalogReward {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	code, rewards := ps.service.TakeReward(p, msg.Type, msg.Id)
	return &protos.ResCatalogReward{
		Code:  int32(code),
		RewardVos: rewards,
	}
}
 
