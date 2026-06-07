package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/catalog"
	player "github.com/forfun/gforgame/internal/service/player"
)

type CatalogRoute struct {
	service *catalog.CatalogService
	player  *player.PlayerService
}

func NewCatalogRoute(service *catalog.CatalogService, playerService *player.PlayerService) *CatalogRoute {
	return &CatalogRoute{
		service: service,
		player:  playerService,
	}
}


func (ps *CatalogRoute) ReqCatalogReward(playerId string, index int32, msg *protos.ReqCatalogReward) *protos.ResCatalogReward {
	p := ps.player.GetPlayer(playerId)
	code, rewards := ps.service.TakeReward(p, msg.Type, msg.Id)
	return &protos.ResCatalogReward{
		Code:  int32(code),
		RewardVos: rewards,
	}
}
 
