package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/monthcard"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
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
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	err := ps.service.TakeReward(player, msg.Type)
	if err != nil {
		return &protos.ResMonthCardGetReward{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMonthCardGetReward{}
}