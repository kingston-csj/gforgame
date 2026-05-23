package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/monthcard"
	player "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
)

type MonthCardRoute struct {
	network.Base
	service *monthcard.MonthCardService
	player  *player.PlayerService
}

func NewMonthCardRoute(service *monthcard.MonthCardService, playerService *player.PlayerService) *MonthCardRoute {
	return &MonthCardRoute{
		service: service,
		player:  playerService,
	}
}

func (ps *MonthCardRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
	context.EventBus.Subscribe(events.Recharge, func(data interface{}) {
		evt := data.(*events.RechargeEvent)
		player := evt.Player.(*playerdomain.Player)
		ps.service.OnRecharge(player, evt.RechargeId)
	})
}


func (ps *MonthCardRoute) ReqGetReward(s *network.Session, index int32, msg *protos.ReqMonthCardGetReward) *protos.ResMonthCardGetReward{
	player := ps.player.GetPlayerBySession(s)
	err := ps.service.TakeReward(player, msg.Type)
	if err != nil {
		return &protos.ResMonthCardGetReward{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMonthCardGetReward{}
}