package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/mall"
	player "github.com/forfun/gforgame/internal/service/player"
)

type MallRoute struct {
	service *mall.MallService
	player  *player.PlayerService
}

func NewMallRoute(service *mall.MallService, playerService *player.PlayerService) *MallRoute {
	return &MallRoute{
		service: service,
		player:  playerService,
	}
}


func (ps *MallRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *MallRoute) ReqMallBuy(playerId string, index int32, msg *protos.ReqMallBuy) *protos.ResMallBuy{
	player := ps.player.GetPlayer(playerId)
	err := ps.service.Buy(player, msg.ProductId, msg.Count)
	if err != nil {
		return &protos.ResMallBuy{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMallBuy{}
}