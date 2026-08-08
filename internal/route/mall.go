package route

import (
	"github.com/forfun/gforgame/common/eventbus"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/mall"
)

type MallRoute struct {
	service    *mall.MallService
	playerrepo *playerrepo.PlayerRepository
}

func NewMallRoute(service *mall.MallService, playerRepo *playerrepo.PlayerRepository) *MallRoute {
	return &MallRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (ps *MallRoute) Init() {
	eventbus.Default().Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *MallRoute) ReqMallBuy(playerId string, index int32, msg *protos.ReqMallBuy) *protos.ResMallBuy {
	player := ps.playerrepo.GetPlayer(playerId)
	err := ps.service.Buy(player, msg.ProductId, msg.Count)
	if err != nil {
		return &protos.ResMallBuy{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResMallBuy{}
}
