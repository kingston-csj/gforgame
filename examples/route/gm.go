package route

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/service/gm"
	"strings"

	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type GmRoute struct {
	network.Base
	service *gm.GmService
}

func NewGmRoute() *GmRoute {
	return &GmRoute{
		service: gm.GetGmService(),
	}
}

func (ps *GmRoute) Init() {
}

func (ps *GmRoute) ReqAction(s *network.Session, index int32, msg *protos.ReqGmCommand) interface{} {
	topic := strings.Split(msg.Args, " ")[0]
	params := ""
	if len(strings.Split(msg.Args, " "))>=2 {
		params = strings.Split(msg.Args, " ")[1]
	} 
	player := network.GetPlayerBySession(s)
	err := ps.service.Dispatch(player.(*playerdomain.Player), topic, params)
	if err != nil {
		return &protos.ResGmCommand{Code: int32(err.Code())}
	}

	return &protos.ResGmCommand{Code: 0}
}
