package route

import (
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/monthcard"
)

type MonthCardRoute struct {
	service    *monthcard.MonthCardService
	playerrepo *playerrepo.PlayerRepository
}

func NewMonthCardRoute(service *monthcard.MonthCardService, playerRepo *playerrepo.PlayerRepository) *MonthCardRoute {
	return &MonthCardRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (ps *MonthCardRoute) ReqGetReward(playerId string, index int32, msg *protos.ReqMonthCardGetReward) *protos.ResMonthCardGetReward {
	player := ps.playerrepo.GetPlayer(playerId)
	err := ps.service.TakeReward(player, msg.Type)
	if err != nil {
		return &protos.ResMonthCardGetReward{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMonthCardGetReward{}
}
