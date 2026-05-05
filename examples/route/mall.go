package route

import (
	"github.com/forfun/gforgame/examples/context"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/events"
	"github.com/forfun/gforgame/examples/protos"
	"github.com/forfun/gforgame/examples/service/mall"
	playerservice "github.com/forfun/gforgame/examples/service/player"
	"github.com/forfun/gforgame/network"
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
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
	err := ps.service.Buy(player, msg.ProductId, msg.Count)
	if err != nil {
		return &protos.ResMallBuy{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMallBuy{}
}