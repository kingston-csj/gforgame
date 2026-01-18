package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/mall"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type MallRoute struct {
	network.Base
	service *mall.MallService
}

func NewMallRoute() *MallRoute {
	return &MallRoute{
		service: mall.GetMallService(),
	}
}

func (ps *MallRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *MallRoute) ReqMallBuy(s *network.Session, index int32, msg *protos.ReqMallBuy) *protos.ResMallBuy{
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	err := ps.service.Buy(player, msg.ProductId, msg.Count)
	if err != nil {
		return &protos.ResMallBuy{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMallBuy{}
}