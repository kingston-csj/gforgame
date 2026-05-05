package route

import (
	"github.com/forfun/gforgame/examples/context"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/events"
	"github.com/forfun/gforgame/examples/protos"
	"github.com/forfun/gforgame/examples/service/monthcard"
	playerservice "github.com/forfun/gforgame/examples/service/player"
	"github.com/forfun/gforgame/network"
)

type MonthCardRoute struct {
	network.Base
	service *monthcard.MonthCardService
}

func NewMonthCardRoute() *MonthCardRoute {
	return &MonthCardRoute{
		service: monthcard.GetMonthCardService(),
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
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
	err := ps.service.TakeReward(player, msg.Type)
	if err != nil {
		return &protos.ResMonthCardGetReward{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMonthCardGetReward{}
}