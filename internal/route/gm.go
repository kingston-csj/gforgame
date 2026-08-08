package route

import (
	"strings"

	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/service/gm"

	"github.com/forfun/gforgame/internal/protos"
)

type GmRoute struct {
	service    *gm.GmService
	playerrepo *playerrepo.PlayerRepository
}

func NewGmRoute(service *gm.GmService, playerRepo *playerrepo.PlayerRepository) *GmRoute {
	return &GmRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (ps *GmRoute) ReqAction(playerId string, index int32, msg *protos.ReqGmCommand) interface{} {
	topic := strings.Split(msg.Args, " ")[0]
	params := ""
	if len(strings.Split(msg.Args, " ")) >= 2 {
		params = strings.Split(msg.Args, " ")[1]
	}
	player := ps.playerrepo.GetPlayer(playerId)
	err := ps.service.Dispatch(player, topic, params)
	if err != nil {
		return &protos.ResGmCommand{Code: int32(err.Code())}
	}

	return &protos.ResGmCommand{Code: 0}
}
