package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/monthcard"
	player "github.com/forfun/gforgame/internal/service/player"
)

type MonthCardRoute struct {
	service *monthcard.MonthCardService
	player  *player.PlayerService
}

func NewMonthCardRoute(service *monthcard.MonthCardService, playerService *player.PlayerService) *MonthCardRoute {
	return &MonthCardRoute{
		service: service,
		player:  playerService,
	}
}

func (ps *MonthCardRoute) ReqGetReward(playerId string, index int32, msg *protos.ReqMonthCardGetReward) *protos.ResMonthCardGetReward{
	player := ps.player.GetPlayer(playerId)
	err := ps.service.TakeReward(player, msg.Type)
	if err != nil {
		return &protos.ResMonthCardGetReward{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMonthCardGetReward{}
}