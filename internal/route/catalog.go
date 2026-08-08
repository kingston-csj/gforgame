package route

import (
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/catalog"
)

type CatalogRoute struct {
	service    *catalog.CatalogService
	playerrepo *playerrepo.PlayerRepository
}

func NewCatalogRoute(service *catalog.CatalogService, playerRepo *playerrepo.PlayerRepository) *CatalogRoute {
	return &CatalogRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (ps *CatalogRoute) ReqCatalogReward(playerId string, index int32, msg *protos.ReqCatalogReward) *protos.ResCatalogReward {
	p := ps.playerrepo.GetPlayer(playerId)
	code, rewards := ps.service.TakeReward(p, msg.Type, msg.Id)
	return &protos.ResCatalogReward{
		Code:      int32(code),
		RewardVos: rewards,
	}
}
