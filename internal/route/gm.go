package route

import (
	"strings"

	"github.com/forfun/gforgame/internal/service/gm"
	player "github.com/forfun/gforgame/internal/service/player"

	"github.com/forfun/gforgame/internal/protos"
)

type GmRoute struct {
	service *gm.GmService
	player  *player.PlayerService
}

func NewGmRoute(service *gm.GmService, playerService *player.PlayerService) *GmRoute {
	return &GmRoute{
		service: service,
		player:  playerService,
	}
}

func (ps *GmRoute) ReqAction(playerId string, index int32, msg *protos.ReqGmCommand) interface{} {
	topic := strings.Split(msg.Args, " ")[0]
	params := ""
	if len(strings.Split(msg.Args, " "))>=2 {
		params = strings.Split(msg.Args, " ")[1]
	} 
	player := ps.player.GetPlayer(playerId)
	err := ps.service.Dispatch(player, topic, params)
	if err != nil {
		return &protos.ResGmCommand{Code: int32(err.Code())}
	}

	return &protos.ResGmCommand{Code: 0}
}
