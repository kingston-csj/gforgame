package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/catalog"
	player "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
)

type CatalogRoute struct {
	network.Base
	service *catalog.CatalogService
	player  *player.PlayerService
}

func NewCatalogRoute(service *catalog.CatalogService, playerService *player.PlayerService) *CatalogRoute {
	return &CatalogRoute{
		service: service,
		player:  playerService,
	}
}


func (ps *CatalogRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *CatalogRoute) ReqCatalogReward(s *network.Session, index int32, msg *protos.ReqCatalogReward) *protos.ResCatalogReward {
	p := ps.player.GetPlayerBySession(s)
	code, rewards := ps.service.TakeReward(p, msg.Type, msg.Id)
	return &protos.ResCatalogReward{
		Code:  int32(code),
		RewardVos: rewards,
	}
}
 
