package route

import (
	"github.com/forfun/gforgame/examples/context"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/events"
	"github.com/forfun/gforgame/examples/protos"
	"github.com/forfun/gforgame/examples/service/catalog"
	playerservice "github.com/forfun/gforgame/examples/service/player"
	"github.com/forfun/gforgame/network"
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
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	code, rewards := ps.service.TakeReward(p, msg.Type, msg.Id)
	return &protos.ResCatalogReward{
		Code:  int32(code),
		RewardVos: rewards,
	}
}
 
